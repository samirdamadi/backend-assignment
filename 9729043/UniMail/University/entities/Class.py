from django.db import models


class Class(models.Model):
    name = models.CharField(max_length=100, null=False, blank=False)
    teacher = models.CharField(max_length=100, null=False, blank=False)
