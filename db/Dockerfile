# Dockerfile for PostgreSQL Database
FROM postgres:17

# Set environment variables for PostgreSQL configuration
ENV POSTGRES_DB=eclypsiumdb
ENV POSTGRES_USER=applicant
ENV POSTGRES_PASSWORD=goodluck

# Copy initialization SQL script
COPY ./01-database.sql /docker-entrypoint-initdb.d/01-database.sql
COPY ./02-sample-data.sql /docker-entrypoint-initdb.d/02-sample-data.sql

# Optional: Add any additional configuration or setup scripts
# COPY pg_hba.conf /etc/postgresql/pg_hba.conf
# COPY postgresql.conf /etc/postgresql/postgresql.conf

# Expose PostgreSQL default port
EXPOSE 5432
