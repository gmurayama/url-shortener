BEGIN;

CREATE TABLE IF NOT EXISTS public.urls
(
    shortened character varying(7) NOT NULL,
    url character varying(32000) NOT NULL,
    PRIMARY KEY (shortened)
);
END;
