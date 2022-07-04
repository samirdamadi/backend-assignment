from django.core.validators import MaxValueValidator, MinValueValidator
from django.db import models

from University.entities.Class import Class


class Student(models.Model):
    name = models.CharField(max_length=100, null=False, blank=False)
    email = models.CharField(max_length=100, null=False, blank=False)
    class_id = models.ForeignKey(to=Class, on_delete=models.deletion.SET_NULL, null=True)
    score = models.FloatField(validators=[MaxValueValidator(20), MinValueValidator(0)])
