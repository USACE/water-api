FROM postgis/postgis:13-3.1

RUN apt-get update -y
RUN apt-get install wget unzip postgis -y
RUN apt-get clean
RUN rm -rf /var/cache/apt/lists

RUN mkdir /gisdata
COPY nation_script_load.sh .

ENTRYPOINT ["sh", "nation_script_load.sh"]