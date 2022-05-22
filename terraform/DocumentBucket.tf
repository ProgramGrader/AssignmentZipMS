

resource "aws_s3_bucket" "url_s3_doc" {
  bucket = "assignment_doc_bucket"

  force_destroy = true

  tags = {
    Name = "url bucket"
    Environment = "dev"
  }
  #  logging {
  #    target_bucket = "target-bucket"
  #  }
}

resource "aws_s3_object" "install_doc" {
  bucket = aws_s3_bucket.url_s3_doc.bucket_domain_name
  key    = "How_to_Install_Jetbrains_Toolbox_and_IDEs"
  source = "..\\src\\doc_files\\How_to_Install_Jetbrains_Toolbox_and_IDEs.pdf"
}

