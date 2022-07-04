package com.emailProject.emailProject.service;

import com.emailProject.emailProject.model.Course;
import com.emailProject.emailProject.repository.CourseRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.List;

@Service
public class CSVService {
    @Autowired
    CourseRepository repository;
    public void save(MultipartFile file) {
        try {
            List<Course> tutorials = CSVHelper.csvToTutorials(file.getInputStream());
            repository.saveAll(tutorials);
        } catch (IOException e) {
            throw new RuntimeException("fail to store csv data: " + e.getMessage());
        }
    }
    public List<Course> getAllTutorials() {
        return repository.findAll();
    }
}