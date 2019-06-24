-- create type entry_sorting_direction as enum('asc', 'desc');
alter table users add column entry_direction text default 'asc';
