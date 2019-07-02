CREATE TABLE `users` (
  `id` INT(10),
  -- Based on firebase User uid. Not well documented, but if you dig INTo this:
  -- https://firebase.google.com/docs/auth/admin/manage-users
  -- we find that an assigned ID may be up to 128 characters, though generated
  -- uids are 28 characters at time of writing (2019-03-08), though that's not
  -- guaranteed
  `auth_id` VARCHAR(128),
  `legal_id` VARCHAR(128),
  `legal_id_type` ENUM('SSN', 'EIN'),
  `active` TINYINT(1) DEFAULT 1 NOT NULL,
  CONSTRAINT `users_key` PRIMARY KEY ( `id` ),
  CONSTRAINT `users_auth_id_unique` UNIQUE (`auth_id`),
  CONSTRAINT `users_ref_entities` FOREIGN KEY ( `id` ) REFERENCES `entities` ( `id` )
);

CREATE TABLE `apps_users` (
  `app_id` INT(10) NOT NULL,
  `user_id` INT(10) NOT NULL,
  `default_context` INT (10),
  CONSTRAINT `apps_users_key` PRIMARY KEY ( `app_id`, `user_id` ),
  CONSTRAINT `apps_users_default_context_ref_entities` FOREIGN KEY ( `default_context` ) REFERENCES `entities` ( `id` )
);
