-- PostgreSQL 数据库完整DDL脚本
-- 包含表创建和索引优化

-- ========================================
-- 表创建 DDL
-- ========================================

-- Bot 表
CREATE TABLE bot (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    code TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ(6)
);

-- ========================================
-- Bot 表索引
-- ========================================
-- 为 user_id 创建索引，用于查询用户的机器人列表
CREATE INDEX idx_bot_user_id ON bot (user_id);
CREATE INDEX idx_bot_title ON bot (title);
CREATE INDEX idx_bot_created_at ON bot (created_at);
CREATE INDEX idx_bot_user_created_at ON bot (user_id, created_at DESC);

-- Record 表
CREATE TABLE record (
    id SERIAL PRIMARY KEY,
    a_id INTEGER NOT NULL,
    a_sx INTEGER NOT NULL,
    a_sy INTEGER NOT NULL,
    b_id INTEGER NOT NULL,
    b_sx INTEGER NOT NULL,
    b_sy INTEGER NOT NULL,
    a_steps TEXT NOT NULL,
    b_steps TEXT NOT NULL,
    game_map TEXT NOT NULL,
    loser_name VARCHAR(255) Not NULL,
    loser_id   BIGINT not NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ(6)
);

-- ========================================
-- Record 表索引
-- ========================================

CREATE INDEX idx_record_created_at ON record (created_at DESC);
CREATE INDEX idx_record_a_id_created_at ON record (a_id, created_at DESC);
CREATE INDEX idx_record_b_id_created_at ON record (b_id, created_at DESC);
CREATE INDEX idx_record_players ON record (a_id, b_id);
CREATE INDEX idx_record_loser ON record (loser_id);

-- User 表
CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) not NULL,
    password VARCHAR(255) not NULL,
    photo VARCHAR(500),
    rating INTEGER not NULL DEFAULT 1500,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ(6)
);

-- ========================================
-- User 表索引
-- ========================================
CREATE UNIQUE INDEX idx_user_username ON "user" (username);
CREATE INDEX idx_user_rating ON "user" (rating DESC);
CREATE INDEX idx_user_rating_username ON "user" (rating DESC, username);