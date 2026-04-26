# Delta for security-core-jwt-validation

## MODIFIED Requirements

### Requirement: Claims model behavior and compatibility

The core MUST support standard claims (`iss`, `sub`, `aud`, `exp`, `nbf`, `iat`, `jti`) and MAY include private claims. It SHALL accept `aud` as string or array and normalize to one form. The `Claims` struct MUST expose `Email`, `Name`, `Picture` as `string` and `Roles` as `[]string` typed fields populated from the standard OIDC JWT keys `email`, `name`, `picture`, and `roles`. Non-standard claims not covered by typed fields SHALL remain accessible via the `Private` map. Missing optional typed fields SHALL default to zero values (`""` / `nil`) without failing parsing.
(Previously: Claims model only referenced standard RFC 7519 claims; OIDC profile fields landed in Private map.)

#### Scenario: Audience encodings normalize consistently

- GIVEN equivalent payloads with `aud` as string/array
- WHEN claims are parsed
- THEN both produce equivalent normalized claims
- AND missing optional claims do not fail parsing by themselves

#### Scenario: Typed fields populated from OIDC JWT claims

- GIVEN a JWT with claims `email`, `name`, `picture`, and `roles` in its payload
- WHEN the token is validated and claims are returned
- THEN `claims.Email`, `claims.Name`, `claims.Picture` equal the respective string values
- AND `claims.Roles` equals the roles slice from the token

#### Scenario: Mixed standard and custom claims coexist

- GIVEN a JWT containing typed OIDC fields (`email`, `name`) and an additional non-standard claim (`tenant_id`)
- WHEN the token is validated
- THEN typed fields are populated correctly
- AND `Private["tenant_id"]` holds the non-standard value
- AND no typed field is overwritten by Private map iteration

#### Scenario: Missing optional typed fields default to zero values

- GIVEN a valid JWT with no `email`, `name`, `picture`, or `roles` claims
- WHEN the token is validated
- THEN `claims.Email`, `claims.Name`, `claims.Picture` are empty strings
- AND `claims.Roles` is nil
- AND validation does not fail due to missing optional fields

#### Scenario: Only bare `roles` key is mapped to typed field

- GIVEN a JWT containing a namespaced claim key (e.g. `https://example.com/roles`) but no bare `roles` key
- WHEN the token is validated
- THEN `claims.Roles` is nil
- AND the namespaced key is accessible via `Private`
