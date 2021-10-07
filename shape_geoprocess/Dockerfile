FROM python:3.9

ENV PYTHONUNBUFFERED=1 GDAL_DISABLE_READDIR_ON_OPEN=EMPTY_DIR

RUN useradd -ms /bin/bash appuser

WORKDIR /usr/src/app

USER appuser

COPY ./app/ ./

RUN pip install --no-cache-dir -r requirements.txt

# CMD [ "python", "./run_container.py" ]
CMD [ "python", "./processor.py" ]
