-- +goose Up
ALTER TABLE "users" ADD COLUMN "hashed_password" varchar(256) DEFAULT 'unset' NOT NULL;--> statement-breakpoint
ALTER TABLE "users" ADD CONSTRAINT "users_hashed_password_unique" UNIQUE("hashed_password");

-- +goose Down
ALTER TABLE "users" DROP CONSTRAINT "users_hashed_password_unique";
ALTER TABLE "users" DROP COLUMN "hashed_password";
