create
extension pg_trgm;
/* trigram extension for faster text-search */

drop table if exists val_info;
create table val_info
(
    moniker             text             not null,
    operator_addr       text             not null,
    operator_addr_hex   text             not null,
    self_stake_addr     text             not null,
    reward_addr         text,
    jailed              boolean          not null,
    status              int              not null,
    voting_power        text             not null,
    identity            text,
    website             text,
    details             text,
    security_contact    text,
    commission_rates    double precision not null,
    max_rate            double precision not null,
    max_change_rate     double precision not null,
    min_self_delegation text,
    primary key (operator_addr)
);

drop table if exists val_stats;
create table val_stats
(
    moniker         text,
    operator_addr   text   not null,
    start_block     bigint not null default 0,
    end_block       bigint not null default 0,
    sign_num        int    not null default 0,
    missed_sign_num int    not null default 0,
    proposals_num   int    not null default 0,
    uptime          float  not null default 0,
    primary key (operator_addr, start_block, end_block)
);

drop table if exists val_sign_p;
create table val_sign_p
(
    moniker       text    not null,
    operator_addr text    not null,
    block_height  bigint  not null default 0,
    status        int     not null, /* Can be 0 = scheduled, 1 executed, 2 missed */
    double_sign   boolean not null,
    child_table   int     not null,
    primary key (operator_addr, block_height, child_table)
) PARTITION BY LIST (child_table);

CREATE TABLE val_sign_0 PARTITION OF val_sign_p FOR VALUES IN
(
    0
);
CREATE TABLE val_sign_1 PARTITION OF val_sign_p FOR VALUES IN
(
    1
);
CREATE TABLE val_sign_2 PARTITION OF val_sign_p FOR VALUES IN
(
    2
);
CREATE TABLE val_sign_3 PARTITION OF val_sign_p FOR VALUES IN
(
    3
);
CREATE TABLE val_sign_4 PARTITION OF val_sign_p FOR VALUES IN
(
    4
);
CREATE TABLE val_sign_5 PARTITION OF val_sign_p FOR VALUES IN
(
    5
);
CREATE TABLE val_sign_6 PARTITION OF val_sign_p FOR VALUES IN
(
    6
);
CREATE TABLE val_sign_7 PARTITION OF val_sign_p FOR VALUES IN
(
    7
);
CREATE TABLE val_sign_8 PARTITION OF val_sign_p FOR VALUES IN
(
    8
);
CREATE TABLE val_sign_9 PARTITION OF val_sign_p FOR VALUES IN
(
    9
);

drop table if exists val_sign_missed;
create table val_sign_missed
(
    moniker       text   not null,
    operator_addr text   not null,
    block_height  bigint not null default 0,
    primary key (operator_addr, block_height)
);

drop table if exists proposal_assignments_p;
create table proposal_assignments_p
(
    moniker       text   not null,
    operator_addr text   not null,
    block_height  bigint not null default 0,
    child_table   int    not null,
    primary key (operator_addr, block_height, child_table)
) PARTITION BY LIST (child_table);

CREATE TABLE proposal_assignments_0 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    0
);
CREATE TABLE proposal_assignments_1 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    1
);
CREATE TABLE proposal_assignments_2 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    2
);
CREATE TABLE proposal_assignments_3 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    3
);
CREATE TABLE proposal_assignments_4 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    4
);
CREATE TABLE proposal_assignments_5 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    5
);
CREATE TABLE proposal_assignments_6 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    6
);
CREATE TABLE proposal_assignments_7 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    7
);
CREATE TABLE proposal_assignments_8 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    8
);
CREATE TABLE proposal_assignments_9 PARTITION OF proposal_assignments_p FOR VALUES IN
(
    9
);

drop table if exists proposal;
create table proposal
(
    proposal_id       int  not null,
    voting_start_time text not null,
    voting_end_time   text not null,
    description       text not null,
    moniker           text not null,
    operator_addr     text not null,
    status            int  not null default 0, /*0 = not voted， 1=yes, 2=no, 3=NoWithVeto, 4 = Abstain*/
    primary key (proposal_id, operator_addr)
);


drop table if exists val_ranking;
create table val_ranking
(
    moniker           text   not null,
    operator_addr     text   not null,
    block_height      bigint not null default 0,
    real_voting_power bigint not null,
    ranking           int    not null default 0,
    primary key (operator_addr, block_height)
);