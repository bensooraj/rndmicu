export AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com > home/ec2-user/error.log
docker pull $AWS_ACCOUNT_ID.dkr.ecr.ap-south-1.amazonaws.com/rndmicu:latest