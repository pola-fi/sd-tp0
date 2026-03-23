FROM python:3.9.7-slim
WORKDIR /app
COPY server /app
RUN python -m unittest tests/test_common.py
ENTRYPOINT ["/bin/sh"]
