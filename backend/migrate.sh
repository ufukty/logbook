#!/usr/bin/env bash

alias psql='psql -U ufuktan -d postgres'


for DB in (logbook_usrs logbook_objectives logbook_)
psql -c "DROP DATABASE IF EXISTS; " -c "CREATE DATABASE "