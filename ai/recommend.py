import os
import boto3
from flask import Flask, jsonify
from botocore.exceptions import BotoCoreError, ClientError

app = Flask(__name__)

profile_name = os.getenv('PY_AWS_PROFILE', 'default')
region_name = os.getenv('PY_AWS_REGION', 'us-east-1')
recommender_arn = os.getenv('PY_RECOMMENDER_ARN')

session = boto3.Session(profile_name=profile_name, region_name=region_name)

# Personalizeのランタイムクライアントを作成
personalize_runtime = session.client('personalize-runtime')

@app.route('/')
def hello_world():
    return 'Hello, World!!'

@app.route('/recommendations/<user_id>', methods=['GET'])
def get_recommendations(user_id):
    try:
        # GetRecommendations APIを呼び出してレコメンデーションを取得
        response = personalize_runtime.get_recommendations(
            recommenderArn=recommender_arn,
            userId=user_id,
            numResults=5
        )
        
        # レコメンデーションの結果をJSONで返す
        recommendations = [{"itemId": item["itemId"]} for item in response['itemList']]
        return jsonify(recommendations)
    
    except (BotoCoreError, ClientError) as error:
        return jsonify({"error": str(error)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
