package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/matiasmartin-labs/common-fwk/security/claims"
)

// Validator validates JWT values against policy options.
type Validator interface {
	Validate(ctx context.Context, raw string) (claims.Claims, error)
}

type validator struct {
	options Options
	allowed map[string]struct{}
}

// NewValidator constructs a framework-agnostic JWT validator.
func NewValidator(options Options) (Validator, error) {
	resolved, err := options.withDefaults()
	if err != nil {
		return nil, err
	}

	allowed := make(map[string]struct{}, len(resolved.Methods))
	for _, method := range resolved.Methods {
		allowed[method] = struct{}{}
	}

	return &validator{options: resolved, allowed: allowed}, nil
}

func (v *validator) Validate(ctx context.Context, raw string) (claims.Claims, error) {
	if raw == "" {
		return claims.Claims{}, wrap("parse", ErrMalformedToken, errors.New("token is empty"))
	}

	unverified, _, err := jwtlib.NewParser().ParseUnverified(raw, jwtlib.MapClaims{})
	if err != nil {
		return claims.Claims{}, wrap("parse", ErrMalformedToken, err)
	}

	alg := ""
	if unverified.Method != nil {
		alg = unverified.Method.Alg()
	}
	if _, ok := v.allowed[alg]; !ok {
		return claims.Claims{}, wrap("method", ErrInvalidMethod, fmt.Errorf("alg %q not allowed", alg))
	}

	kid, err := kidFromHeader(unverified.Header)
	if err != nil {
		return claims.Claims{}, wrap("parse", ErrMalformedToken, err)
	}

	key, err := v.options.Resolver.Resolve(ctx, kid)
	if err != nil {
		return claims.Claims{}, wrap("resolve-key", ErrKeyResolution, err)
	}

	parser := jwtlib.NewParser(
		jwtlib.WithValidMethods(v.options.Methods),
		jwtlib.WithoutClaimsValidation(),
	)

	verified, err := parser.Parse(raw, func(token *jwtlib.Token) (any, error) {
		if token.Method == nil {
			return nil, errors.New("token method missing")
		}

		if key.Method != "" && token.Method.Alg() != key.Method {
			return nil, fmt.Errorf("token alg %q does not match key method %q", token.Method.Alg(), key.Method)
		}

		return key.Verify, nil
	})
	if err != nil {
		if errors.Is(err, jwtlib.ErrTokenSignatureInvalid) {
			return claims.Claims{}, wrap("verify-signature", ErrInvalidSignature, err)
		}

		return claims.Claims{}, wrap("parse", ErrMalformedToken, err)
	}

	mapped, err := claimsFromToken(verified)
	if err != nil {
		return claims.Claims{}, wrap("parse", ErrMalformedToken, err)
	}

	if v.options.Issuer != "" && mapped.Issuer != v.options.Issuer {
		return claims.Claims{}, wrap("claims", ErrInvalidIssuer, fmt.Errorf("expected %q, got %q", v.options.Issuer, mapped.Issuer))
	}

	if len(v.options.Audience) > 0 && !hasAnyAudience(mapped, v.options.Audience) {
		return claims.Claims{}, wrap("claims", ErrInvalidAudience, fmt.Errorf("expected one of %v, got %v", v.options.Audience, mapped.NormalizedAudience()))
	}

	now := v.options.Now()
	if mapped.ExpiresAt != nil && mapped.ExpiresAt.Before(now) {
		return claims.Claims{}, wrap("claims", ErrExpiredToken, fmt.Errorf("exp %s before now %s", mapped.ExpiresAt.UTC().Format(time.RFC3339), now.UTC().Format(time.RFC3339)))
	}

	if mapped.NotBefore != nil && mapped.NotBefore.After(now) {
		return claims.Claims{}, wrap("claims", ErrNotYetValidToken, fmt.Errorf("nbf %s after now %s", mapped.NotBefore.UTC().Format(time.RFC3339), now.UTC().Format(time.RFC3339)))
	}

	return mapped, nil
}

func kidFromHeader(header map[string]any) (string, error) {
	v, ok := header["kid"]
	if !ok {
		return "", nil
	}

	kid, ok := v.(string)
	if !ok {
		return "", errors.New("kid header must be string")
	}

	return kid, nil
}

func claimsFromToken(token *jwtlib.Token) (claims.Claims, error) {
	mapClaims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok {
		return claims.Claims{}, errors.New("unexpected claims type")
	}

	mapped := claims.Claims{}
	if issuer, ok := mapClaims["iss"].(string); ok {
		mapped.Issuer = issuer
	}
	if subject, ok := mapClaims["sub"].(string); ok {
		mapped.Subject = subject
	}
	if jwtID, ok := mapClaims["jti"].(string); ok {
		mapped.JWTID = jwtID
	}

	audience, err := parseAudience(mapClaims["aud"])
	if err != nil {
		return claims.Claims{}, err
	}
	mapped.Audience = audience

	if exp, ok, err := parseNumericDate(mapClaims["exp"]); err != nil {
		return claims.Claims{}, fmt.Errorf("exp: %w", err)
	} else if ok {
		mapped.ExpiresAt = &exp
	}

	if nbf, ok, err := parseNumericDate(mapClaims["nbf"]); err != nil {
		return claims.Claims{}, fmt.Errorf("nbf: %w", err)
	} else if ok {
		mapped.NotBefore = &nbf
	}

	if iat, ok, err := parseNumericDate(mapClaims["iat"]); err != nil {
		return claims.Claims{}, fmt.Errorf("iat: %w", err)
	} else if ok {
		mapped.IssuedAt = &iat
	}

	private := make(map[string]interface{})
	for key, value := range mapClaims {
		switch key {
		case "iss", "sub", "aud", "exp", "nbf", "iat", "jti":
			continue
		default:
			private[key] = value
		}
	}
	if len(private) > 0 {
		mapped.Private = private
	}

	return mapped, nil
}

func parseAudience(raw any) (claims.Audience, error) {
	if raw == nil {
		return nil, nil
	}

	switch value := raw.(type) {
	case string:
		return claims.Audience{value}, nil
	case []string:
		out := make([]string, len(value))
		copy(out, value)
		return claims.Audience(out), nil
	case []any:
		out := make([]string, 0, len(value))
		for _, elem := range value {
			s, ok := elem.(string)
			if !ok {
				return nil, claims.ErrInvalidAudienceType
			}
			out = append(out, s)
		}
		return claims.Audience(out), nil
	default:
		return nil, claims.ErrInvalidAudienceType
	}
}

func parseNumericDate(raw any) (time.Time, bool, error) {
	if raw == nil {
		return time.Time{}, false, nil
	}

	var seconds int64
	switch value := raw.(type) {
	case float64:
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return time.Time{}, false, errors.New("invalid numeric date value")
		}
		seconds = int64(value)
	case float32:
		fv := float64(value)
		if math.IsNaN(fv) || math.IsInf(fv, 0) {
			return time.Time{}, false, errors.New("invalid numeric date value")
		}
		seconds = int64(value)
	case int64:
		seconds = value
	case int32:
		seconds = int64(value)
	case int:
		seconds = int64(value)
	case json.Number:
		parsed, err := strconv.ParseInt(value.String(), 10, 64)
		if err != nil {
			return time.Time{}, false, err
		}
		seconds = parsed
	case string:
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return time.Time{}, false, err
		}
		seconds = parsed
	default:
		return time.Time{}, false, fmt.Errorf("unsupported numeric date type %T", raw)
	}

	return time.Unix(seconds, 0).UTC(), true, nil
}

func hasAnyAudience(claimsValue claims.Claims, expected []string) bool {
	for _, candidate := range expected {
		if claimsValue.HasAudience(candidate) {
			return true
		}
	}

	return false
}
