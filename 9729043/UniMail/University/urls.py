from django.urls import path

from .views import ClassView, CreateClassView, StudentView, CreateStudentView, EmailView

urlpatterns = [
    path('classes/<int:class_id>', ClassView.as_view(), name='single_class'),
    path('classes', CreateClassView.as_view(), name='create_class'),
    path('students/<int:student_id>', StudentView.as_view(), name='single_student'),
    path('students', CreateStudentView.as_view(), name='create_student'),
    path('mail/send', EmailView.as_view(), name='send_email')
]
