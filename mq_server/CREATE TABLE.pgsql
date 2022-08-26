DROP TABLE if exists public.client cascade
DROP TABLE public.transaction

CREATE TABLE public.client(
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 name VARCHAR(100) UNIQUE NOT NULL,
 ip VARCHAR(20) NOT NULL,
 balance INTEGER
);

CREATE TABLE public.transaction(
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 info VARCHAR(100) NOT NULL,
 sum INTEGER,
 client_id   UUID NOT NULL,
CONSTRAINT client_fk FOREIGN KEY (client_id) REFERENCES public.client(id)
);


INSERT INTO transaction (info,sum,client_id) VALUES ('Add 200 rub',200,'c46ac469-7b98-4818-8a27-6b37c88d7628');

INSERT INTO client (name,ip,balance) VALUES ('testClient','127.0.0.1:9999',0);
