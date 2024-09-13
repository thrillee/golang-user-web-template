data "external_schema" "gorm" {
   program = ["go", "run", "main.go", "makemigrations"]
}

variable "url" {
  type = string
  default = getenv("DATABASE_URL")
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = var.url
  url = var.url
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
