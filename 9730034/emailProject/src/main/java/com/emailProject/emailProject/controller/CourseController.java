package com.emailProject.emailProject.controller;

import com.emailProject.emailProject.model.ResponseMessage;
import com.emailProject.emailProject.service.AuthorizationService;
import com.emailProject.emailProject.model.Course;
import com.emailProject.emailProject.model.Student;
import com.emailProject.emailProject.repository.CourseRepository;
import com.emailProject.emailProject.repository.StudentRepository;
import com.emailProject.emailProject.service.CSVHelper;
import com.emailProject.emailProject.service.CSVService;
import com.emailProject.emailProject.service.MailService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.Optional;

@RequestMapping(path = "/course")
@RestController
public class CourseController {
    @Autowired
    private CourseRepository courseRepository;

    @Autowired
    private StudentRepository studentRepository;

    @Autowired
    private AuthorizationService authorizationService;

    @Autowired
    private MailService mailService;

    @Autowired
    private CSVService csvService;

    @PostMapping(path = "/create")
    public ResponseEntity<Course> createCourse(@RequestBody Course course, @RequestHeader("authorization") String auth){
        if (authorizationService.haveAccess(auth))
            return new ResponseEntity<Course> (courseRepository.save(course), HttpStatus.OK);
        else
            return new ResponseEntity<Course>( HttpStatus.UNAUTHORIZED);
    }

    @GetMapping(path = "/all")
    public ResponseEntity<List<Course>> findAll(@RequestHeader("authorization") String auth){
        if (authorizationService.haveAccess(auth))
            return new ResponseEntity<List<Course>> ( courseRepository.findAll(), HttpStatus.OK);
        else return new ResponseEntity<List<Course>>( HttpStatus.UNAUTHORIZED);
    }

    @DeleteMapping(path = "/delete/{courseId}")
    public ResponseEntity<List<Course>> delete(@PathVariable("courseId") Integer courseId, @RequestHeader("authorization") String auth){
        if (authorizationService.haveAccess(auth)) {
            Optional<Course> course = courseRepository.findById(Long.valueOf(courseId));
            for (Student s : course.get().getStudents()) {
                s.setCourse(null);
                studentRepository.save(s);
            }
            course.get().setStudents(null);
            courseRepository.save(course.get());
            courseRepository.delete(course.get());
            return new ResponseEntity<List<Course>> ( courseRepository.findAll(), HttpStatus.OK);
        }
        else return new ResponseEntity<List<Course>>( HttpStatus.UNAUTHORIZED);
    }

    @PostMapping(path = "/sendmail/{courseId}")
    public ResponseEntity<Map<Student, Boolean>> sendEmail(@PathVariable("courseId") Integer courseId, @RequestHeader("authorization") String auth){
        if (authorizationService.haveAccess(auth)) {
            Optional<Course> course = courseRepository.findById(Long.valueOf(courseId));
            return new ResponseEntity<Map<Student, Boolean>>(mailService.sendMessages(course.get()), HttpStatus.OK);
        } else return new ResponseEntity<Map<Student, Boolean>>( HttpStatus.UNAUTHORIZED);
    }

    @PostMapping("/upload")
    public ResponseEntity<ResponseMessage> uploadFile(@RequestParam("file") MultipartFile file, @RequestHeader("authorization") String auth) {
        if (authorizationService.haveAccess(auth)) {
            String message = "";
            if (CSVHelper.hasCSVFormat(file)) {
                try {
                    csvService.save(file);
                    message = "Uploaded the file successfully: " + file.getOriginalFilename();
                    return ResponseEntity.status(HttpStatus.OK).body(new ResponseMessage(message));
                } catch (Exception e) {
                    message = "Could not upload the file: " + file.getOriginalFilename() + "!";
                    return ResponseEntity.status(HttpStatus.EXPECTATION_FAILED).body(new ResponseMessage(message));
                }
            }
            message = "Please upload a csv file!";
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(new ResponseMessage(message));
        }
        else return ResponseEntity.status( HttpStatus.UNAUTHORIZED).body(new ResponseMessage("un authorized"));
    }
}
