-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
                               id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                               service_name TEXT NOT NULL,
                               price        BIGINT NOT NULL,
                               user_id      UUID NOT NULL,
                               start_date   DATE NOT NULL,
                               end_date     DATE NOT NULL,
                               updated_at   TIMESTAMP WITH TIME ZONE DEFAULT now(),
                               created_at   TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);