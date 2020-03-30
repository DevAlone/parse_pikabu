# parse_pikabu

Backend part of pikagraphs project.

See how it works here -> https://pikastat.d3d.info

# Install

1 Install postgres
2 Create db and user
```sql
postgres> CREATE USER username WITH ENCRYPTED PASSWORD 'password';
postgres> CREATE DATABASE database_name;
postgres> GRANT ALL ON DATABASE database_name TO username;
username> CREATE EXTENSION pg_trgm;
```

# Mirrors

- https://github.com/DevAlone/parse_pikabu
- https://gitlab.com/DevAlone/parse_pikabu
- https://bitbucket.org/d3dev/parse_pikabu
