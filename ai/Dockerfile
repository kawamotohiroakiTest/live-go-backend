FROM python:3.12-slim

RUN apt-get update && apt-get install -y \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# boto3 と awscli をインストール
RUN pip install --no-cache-dir boto3 awscli

# Flaskなどの依存パッケージをインストール
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

WORKDIR /app

COPY . .

# Flaskサーバーを起動
CMD ["python3", "recommend.py"]
