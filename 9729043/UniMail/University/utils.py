from django.core.mail import send_mail

from UniversityMail.settings import EMAIL_HOST_USER
from .models import Student


def send_email(student: Student):
    content = f"Name: {student.name}\r\nClassName: {student.class_id.name}" + \
              f"\r\nTeacher: {student.class_id.teacher}\r\nScore: {student.score}"
    recipient = student.email
    subject = "Student Details"
    return send_mail(subject, content, EMAIL_HOST_USER, [recipient])
