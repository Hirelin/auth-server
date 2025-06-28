-- AUTH QUERIES
-- name: GetUserByEmail :one
SELECT * FROM "User" WHERE email = $1;

-- name: AddUser :one
INSERT INTO "User" (email, name) VALUES ($1, $2)
RETURNING *;

-- name: CreateSession :one
INSERT INTO "Session" (session_token, refresh_token, expires_at, user_id) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteSession :one
DELETE FROM "Session" WHERE session_token = $1
RETURNING session_token;

-- name: GetSessionUser :one
SELECT id, email, name, image FROM "User" WHERE id = (SELECT user_id FROM "Session" WHERE session_token = $1);

-- name: UpdateSession :one
UPDATE "Session" SET session_token = $1, refresh_token = $2, expires_at = $3 WHERE session_token = $4
RETURNING *;

-- name: CreateOAuthAccount :one
INSERT INTO "Account" (user_id, type, provider, provider_account_id, access_token, refresh_token, expires_at, token_type, scope, id_token, session_state)
VALUES ($1, "oauth", $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateOAuthAccount :one
UPDATE "Account"
SET access_token = $2,
    refresh_token = $3,
    expires_at = $4,
    token_type = $5,
    scope = $6,
    id_token = $7,
    session_state = $8
WHERE user_id = $9 AND provider = $10 AND provider_account_id = $1
RETURNING *;

-- name: NewOAuthUserTransaction :one
-- name: NewOAuthUserTransaction :one
WITH new_user AS (
    INSERT INTO "User" (email, name, image, "emailVerified") 
    VALUES ($1, $2, $3, $4)
    RETURNING *
),
new_account AS (
    INSERT INTO "Account" (user_id, type, provider, provider_account_id, access_token, refresh_token, expires_at, token_type, scope, id_token, session_state)
    VALUES ((SELECT id FROM new_user), 'oauth', $5, $6, $7, $8, $9, $10, $11, $12, $13)
    RETURNING *
)
SELECT 
    new_user.id as user_id,
    new_user.email,
    new_user.name,
    new_user.image,
    new_user."emailVerified",
    new_account.id as account_id,
    new_account.provider,
    new_account.provider_account_id
FROM new_user, new_account;

-- name: GetAccountByProviderId :one
SELECT * FROM "Account" WHERE provider = $1 AND provider_account_id = $2;

-- name: GetVerificationToken :one
SELECT * FROM "VerificationToken" WHERE token = $1;

-- name: CreateVerificationToken :one
INSERT INTO "VerificationToken" (identifier, token, expires_at) VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteVerificationToken :one
DELETE FROM "VerificationToken" WHERE identifier = $1 AND token = $2
RETURNING *;

-- name: UpdateUserVerification :one
UPDATE "User" SET "emailVerified" = $2 WHERE email = $1
RETURNING *;
