-- Adminer 4.7.7 PostgreSQL dump

DROP TABLE IF EXISTS "customers";
DROP SEQUENCE IF EXISTS customers_cust_id_seq;
CREATE SEQUENCE customers_cust_id_seq
INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."customers"
(
    "cust_id" integer DEFAULT nextval('customers_cust_id_seq') NOT NULL,
    "account_num" character varying(10),
    "cust_name" character varying(100) NOT NULL,
    "cust_address" text NOT NULL,
    "cust_phone" character varying(20) NOT NULL,
    "cust_email" character varying(64) NOT NULL,
    "cust_password" character varying(64) NOT NULL,
    "cust_pict" character varying(200),
    "is_verified" boolean,
    "channel" character varying(20),
    "card_num" character varying(20),
    "cvv" character varying(3),
    "expired" timestamp,
    "created_at" timestamp,
    "updated_at" timestamp,
    "is_deleted" timestamp,
    CONSTRAINT "customers_account_num" UNIQUE ("account_num"),
    CONSTRAINT "customers_cust_email_key" UNIQUE ("cust_email"),
    CONSTRAINT "customers_cust_phone_key" UNIQUE ("cust_phone"),
    CONSTRAINT "customers_pkey" PRIMARY KEY ("cust_id")
)
WITH (oids = false);


DROP TABLE IF EXISTS "accounts";
DROP SEQUENCE IF EXISTS accounts_account_id_seq;
CREATE SEQUENCE accounts_account_id_seq
INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."accounts"
(
    "account_id" integer DEFAULT nextval('accounts_account_id_seq') NOT NULL,
    "account_num" character varying(10),
    "account_balance" integer,
    "created_at" timestamp,
    CONSTRAINT "accounts_account_num_key" UNIQUE ("account_num"),
    CONSTRAINT "accounts_pkey" PRIMARY KEY ("account_id"),
    CONSTRAINT "account_fk" FOREIGN KEY (account_num) REFERENCES customers(account_num)
    NOT DEFERRABLE
)
    WITH
    (oids = false);


    DROP TABLE IF EXISTS "email_token";
    DROP SEQUENCE IF EXISTS email_token_et_id_seq;
    CREATE SEQUENCE email_token_et_id_seq
    INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

    CREATE TABLE "public"."email_token"
    (
        "et_id" integer DEFAULT nextval('email_token_et_id_seq') NOT NULL,
        "token" text,
        "email" character varying(64),
        CONSTRAINT "email_token_email" UNIQUE ("email"),
        CONSTRAINT "email_token_pkey" PRIMARY KEY ("et_id"),
        CONSTRAINT "email_token_token_key" UNIQUE ("token"),
        CONSTRAINT "email_token_email_fkey" FOREIGN KEY (email) REFERENCES customers(cust_email) ON UPDATE CASCADE ON DELETE CASCADE
        NOT DEFERRABLE
)
        WITH
        (oids = false);


        DROP TABLE IF EXISTS "partners";
        DROP SEQUENCE IF EXISTS partners_partner_id_seq;
        CREATE SEQUENCE partners_partner_id_seq
        INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

        CREATE TABLE "public"."partners"
        (
            "partner_id" integer DEFAULT nextval('partners_partner_id_seq') NOT NULL,
            "client_id" integer,
            "secret" character varying(200),
            CONSTRAINT "partners_client_id_key" UNIQUE ("client_id"),
            CONSTRAINT "partners_pkey" PRIMARY KEY ("partner_id"),
            CONSTRAINT "partners_secret_key" UNIQUE ("secret")
        )
        WITH (oids = false);


        DROP TABLE IF EXISTS "transaction_logs";
        DROP SEQUENCE IF EXISTS transaction_logs_tl_id_seq1;
        CREATE SEQUENCE transaction_logs_tl_id_seq1
        INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

        CREATE TABLE "public"."transaction_logs"
        (
            "tl_id" integer DEFAULT nextval('transaction_logs_tl_id_seq1') NOT NULL,
            "account_num" character varying(10),
            "from_account" character varying(20),
            "dest_account" character varying(20),
            "tran_amount" integer,
            "description" text,
            "created_at" timestamp,
            CONSTRAINT "transaction_logs_pkey" PRIMARY KEY ("tl_id"),
            CONSTRAINT "fk_account" FOREIGN KEY (account_num) REFERENCES customers(account_num)
            NOT DEFERRABLE
)
            WITH
            (oids = false);


            DROP TABLE IF EXISTS "virtual_accounts";
            DROP SEQUENCE IF EXISTS virtual_accounts_va_id_seq;
            CREATE SEQUENCE virtual_accounts_va_id_seq
            INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

            CREATE TABLE "public"."virtual_accounts"
            (
                "va_id" integer DEFAULT nextval('virtual_accounts_va_id_seq') NOT NULL,
                "va_num" character varying(13),
                "account_num" character varying(10),
                "va_balance" integer,
                "va_color" character varying(15),
                "va_label" character varying(100),
                "created_at" timestamp,
                "updated_at" timestamp,
                CONSTRAINT "virtual_accounts_pkey" PRIMARY KEY ("va_id"),
                CONSTRAINT "virtual_accounts_va_num_key" UNIQUE ("va_num"),
                CONSTRAINT "virtual_accounts_account_num_fkey" FOREIGN KEY (account_num) REFERENCES customers(account_num)
                NOT DEFERRABLE
)
                WITH
                (oids = false);

                DROP TABLE IF EXISTS "admins";
                DROP SEQUENCE IF EXISTS admin_id_seq;
                CREATE SEQUENCE admin_id_seq
                INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

                CREATE TABLE "public"."admins"
                (
                    "id" integer DEFAULT nextval('admin_id_seq') NOT NULL,
                    "username" character varying(20) NOT NULL,
                    "password" character varying(64) NOT NULL,
                    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
                    "lastlogin_at" timestamp,
                    CONSTRAINT "admin_unique" UNIQUE ("username"),
                    CONSTRAINT "admins_pkey" PRIMARY KEY ("id")
                )
                WITH (oids = false);

                DROP TABLE IF EXISTS "log_admins";
                DROP SEQUENCE IF EXISTS log_admins_id_seq;
                CREATE SEQUENCE log_admins_id_seq
                INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

                CREATE TABLE "public"."log_admins"
                (
                    "id" integer DEFAULT nextval('log_admins_id_seq') NOT NULL,
                    "username" character varying(20) NOT NULL,
                    "account_num" character varying(10),
                    "action" character varying(64) NOT NULL,
                    "action_time" timestamp,
                    CONSTRAINT "log_admins_pkey" PRIMARY KEY ("id"),
                    CONSTRAINT "fk_log_admins" FOREIGN KEY (username) REFERENCES admins(username)
                    NOT DEFERRABLE
)
                    WITH
                    (oids = false);

-- 2020-07-25 10:06:53.514544+00