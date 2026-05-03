SELECT setval(
    pg_get_serial_sequence('email_verification_tokens', 'id'),
    COALESCE((SELECT MAX(id) FROM email_verification_tokens), 0) + 1,
    false
);

SELECT setval(
    pg_get_serial_sequence('password_reset_tokens', 'id'),
    COALESCE((SELECT MAX(id) FROM password_reset_tokens), 0) + 1,
    false
);
