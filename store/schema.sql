create table users (
    UID SERIAL,
    VKID INT,
    NODE INT DEFAULT 0,
    UserName TEXT DEFAULT '',

    SelectDate INT DEFAULT 0,
    SelectCollection INT DEFAULT 0,
    SelectType INT DEFAULT 0,
    SelectStoreys INT DEFAULT 0,
    SelectProject INT DEFAULT 0,
    Data TEXT DEFAULT ''
);
CREATE INDEX on users(VKID);