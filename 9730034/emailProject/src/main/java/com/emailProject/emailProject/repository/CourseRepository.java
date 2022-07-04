package com.emailProject.emailProject.repository;

import com.emailProject.emailProject.model.Course;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Component;
import org.springframework.stereotype.Repository;
import org.springframework.stereotype.Service;

@Service
public interface CourseRepository extends JpaRepository<Course, Long> {
}
