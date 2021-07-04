export AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
docker pull $AWS_ACCOUNT_ID.dkr.ecr.ap-south-1.amazonaws.com/rndmicu:latest