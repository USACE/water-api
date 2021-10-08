FROM python:3.9

ENV PYTHONUNBUFFERED=1 GDAL_DISABLE_READDIR_ON_OPEN=EMPTY_DIR

RUN useradd -ms /bin/bash appuser

WORKDIR /usr/src/app

COPY ./app/ ./

RUN pip install --no-cache-dir -r requirements.txt

USER appuser

# CMD [ "python", "./run_container.py" ]
CMD [ "python", "./processor.py" ]
