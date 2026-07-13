terraform {
  required_version = ">= 1.0"

  backend "s3" {
    bucket         = "queueflow-terraform-state"
    key            = "infra/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "queueflow-terraform-locks"
    encrypt        = true
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = "queueflow-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["${var.aws_region}a", "${var.aws_region}b", "${var.aws_region}c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true

  tags = {
    Project = "QueueFlow"
  }
}

# EKS Cluster
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = "queueflow-cluster"
  cluster_version = "1.28"

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  eks_managed_node_groups = {
    queueflow_nodes = {
      min_size     = 1
      max_size     = 4
      desired_size = 2

      instance_types = ["t3.medium"]
      capacity_type  = "ON_DEMAND"
    }
  }

  tags = {
    Project = "QueueFlow"
  }
}

# RDS Postgres
resource "aws_db_subnet_group" "queueflow" {
  name       = "queueflow-db-subnet-group"
  subnet_ids = module.vpc.private_subnets

  tags = {
    Project = "QueueFlow"
  }
}

resource "aws_db_instance" "queueflow_postgres" {
  identifier             = "queueflow-postgres"
  engine                 = "postgres"
  engine_version         = "16.3"
  instance_class         = "db.t3.micro"
  allocated_storage      = 20
  storage_type           = "gp3"
  db_subnet_group_name   = aws_db_subnet_group.queueflow.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  db_name  = "queueflow"
  username = "queueflow"
  password = var.db_password

  skip_final_snapshot = true

  tags = {
    Project = "QueueFlow"
  }
}

# MSK Kafka
resource "aws_msk_cluster" "queueflow_kafka" {
  cluster_name           = "queueflow-kafka"
  kafka_version          = "3.5.1"
  number_of_broker_nodes = 3

  broker_node_group_info {
    instance_type   = "kafka.t3.small"
    ebs_volume_size = 100
    client_subnets  = module.vpc.private_subnets
    security_groups = [aws_security_group.msk.id]
  }

  tags = {
    Project = "QueueFlow"
  }
}

# Security Groups
resource "aws_security_group" "rds" {
  name_prefix = "queueflow-rds-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "msk" {
  name_prefix = "queueflow-msk-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 9092
    to_port     = 9092
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# IAM for EKS
resource "aws_iam_role" "eks_role" {
  name = "queueflow-eks-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "eks.amazonaws.com"
        }
      }
    ]
  })
}

# Outputs
output "kubeconfig" {
  value     = module.eks.kubeconfig
  sensitive = true
}

output "postgres_endpoint" {
  value = aws_db_instance.queueflow_postgres.endpoint
}

output "kafka_bootstrap_brokers" {
  value = aws_msk_cluster.queueflow_kafka.bootstrap_brokers
}

output "vpc_id" {
  value = module.vpc.vpc_id
}
