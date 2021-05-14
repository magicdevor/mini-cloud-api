CREATE TABLE "users" (
  "id" varchar(20) PRIMARY KEY,
  "openid" varchar(50) UNIQUE NOT NULL,
  "session_key" varchar(50) NOT NULL,
  "appid" varchar(50) NOT NULL,
  "unionid" varchar(50) NOT NULL DEFAULT 'unionid',
  "openid_from" varchar(50) NOT NULL,
  "appid_from" varchar(50) NOT NULL,
  "unionid_from" varchar(50) NOT NULL DEFAULT 'unionid',
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "users_profile" (
  "id" varchar(20) PRIMARY KEY,
  "user_id" varchar(20) NOT NULL,
  "nickname" varchar(20) NOT NULL DEFAULT 'unknow',
  "avatar_url" varchar(255) NOT NULL,
  "gender" char(1) NOT NULL DEFAULT '0',
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL
);

ALTER TABLE "users_profile" ADD CONSTRAINT "users_profile_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id");