FROM hasura/graphql-engine:v2.12.1
WORKDIR /root 
COPY ./metadata .
RUN curl -L https://github.com/hasura/graphql-engine/raw/stable/cli/get.sh | bash