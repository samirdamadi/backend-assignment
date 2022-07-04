package com.emailProject.emailProject.controller;

import com.emailProject.emailProject.service.AuthorizationService;
import com.emailProject.emailProject.model.Course;
import com.emailProject.emailProject.model.Student;
import com.emailProject.emailProject.repository.CourseRepository;
import com.emailProject.emailProject.repository.StudentRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;

@RequestMapping(path = "/student")
@RestController
public class StudentController {
    @Autowired
    private StudentRepository studentRepository;

    @Autowired
    private CourseRepository courseRepository;

    @Autowired
    private AuthorizationService authorizationService;

    @PostMapping(path = "/create")
    public ResponseEntity<List<Student>> create(@RequestBody List<Student> students, @RequestHeader("authorization") String auth){
        if (authorizationService.haveAccess(auth))
            return  new ResponseEntity<List<Student>>(studentRepository.saveAll(students), HttpStatus.OK);
        else return new ResponseEntity<List<Student>>( HttpStatus.UNAUTHORIZED);
    }

    @GetMapping(path = "/all")
    public ResponseEntity<List<Student>> findAll( @RequestHeader("authorization") String auth){
        if (authorizationService.haveAccess(auth))
            return  new ResponseEntity<List<Student>>(studentRepository.findAll(), HttpStatus.OK);
        else return new ResponseEntity<List<Student>>( HttpStatus.UNAUTHORIZED);
    }

    @DeleteMapping(path = "/delete/{studentId}")
    public ResponseEntity<List<Student>> delete(@PathVariable("studentId") Integer studentId, @RequestHeader("authorization") String auth){
        if (authorizationService.haveAccess(auth)) {
            Optional<Student> student = studentRepository.findById(Long.valueOf(studentId));
            studentRepository.delete(student.get());
            return new ResponseEntity<List<Student>>(studentRepository.findAll(), HttpStatus.OK);

        }
        else  return new ResponseEntity<List<Student>>(HttpStatus.UNAUTHORIZED);

    }

    @PostMapping(value = "/addcourse/{studentId}/{courseId}")
    public ResponseEntity<Optional<Student>> addCourse(@PathVariable("studentId") Integer studentId, @PathVariable("courseId") Integer courseId, @RequestHeader("authorization") String auth){
        if (authorizationService.haveAccess(auth)) {

            Optional<Student> student = studentRepository.findById(Long.valueOf(studentId));
            Optional<Course> course = courseRepository.findById(Long.valueOf(courseId));
            if (student.isEmpty() || course.isEmpty()) {
                return null;
            } else {
                student.get().setCourse(course.get());

                return new ResponseEntity<Optional<Student>>(Optional.of(studentRepository.save(student.get())), HttpStatus.OK);

            }
        }
        else return new ResponseEntity<Optional<Student>>(HttpStatus.UNAUTHORIZED);
    }
}
