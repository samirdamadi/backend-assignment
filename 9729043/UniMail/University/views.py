from rest_framework import status
from rest_framework.exceptions import NotFound, PermissionDenied
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView

from .models import Class, Student
from .serializers import ClassSerializer, StudentSerializer
from .utils import send_email


class CreateClassView(APIView):
    #permission_classes = (IsAuthenticated,)

    def post(self, request):
        serializer = ClassSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        obj = serializer.save()
        response = ClassSerializer(obj).data
        return Response(response, status=status.HTTP_201_CREATED)


class ClassView(APIView):
    #permission_classes = (IsAuthenticated,)

    def get(self, request, class_id):
        try:
            obj = Class.objects.get(pk=class_id)
        except Class.DoesNotExist:
            raise NotFound("class with this id does not exist")
        response = ClassSerializer(obj).data
        return Response(response)

    def delete(self, request, class_id):
        try:
            obj = Class.objects.get(pk=class_id)
        except Class.DoesNotExist:
            raise PermissionDenied()
        obj.delete()
        return Response(status=status.HTTP_204_NO_CONTENT)

    def put(self, request, class_id):
        try:
            obj = Class.objects.get(pk=class_id)
        except Class.DoesNotExist:
            raise PermissionDenied()
        serializer = ClassSerializer(instance=obj, data=request.data)
        serializer.is_valid(raise_exception=True)
        obj = serializer.save()
        response = ClassSerializer(obj).data
        return Response(response)


class CreateStudentView(APIView):
    #permission_classes = (IsAuthenticated,)

    def post(self, request):
        many = type(request.data) is list
        serializer = StudentSerializer(data=request.data, many=many)
        serializer.is_valid(raise_exception=True)
        objs = serializer.save()
        response = StudentSerializer(objs, many=many).data
        return Response(response, status=status.HTTP_201_CREATED)


class StudentView(APIView):
    #permission_classes = (IsAuthenticated,)

    def get(self, request, student_id):
        try:
            obj = Student.objects.get(pk=student_id)
        except Student.DoesNotExist:
            raise NotFound("student with this id does not exist")
        response = StudentSerializer(obj).data
        return Response(response)

    def delete(self, request, student_id):
        try:
            obj = Student.objects.get(pk=student_id)
        except Student.DoesNotExist:
            raise PermissionDenied()
        obj.delete()
        return Response(status=status.HTTP_204_NO_CONTENT)

    def put(self, request, student_id):
        try:
            obj = Student.objects.get(pk=student_id)
        except Student.DoesNotExist:
            raise PermissionDenied()
        serializer = StudentSerializer(instance=obj, data=request.data)
        serializer.is_valid(raise_exception=True)
        obj = serializer.save()
        response = StudentSerializer(obj).data
        return Response(response)


class EmailView(APIView):
    #permission_classes = (IsAuthenticated,)

    def post(self, request):
        class_id = int(self.request.GET.get('class_id'))
        try:
            obj = Class.objects.get(pk=class_id)
        except Class.DoesNotExist:
            raise PermissionDenied()

        students = obj.student_set.all()
        response = list()
        for std in students:
            res = send_email(std)
            response.append({'email': std.email, 'status': 'failed' if res == 0 else 'success'})
        return Response(response)
