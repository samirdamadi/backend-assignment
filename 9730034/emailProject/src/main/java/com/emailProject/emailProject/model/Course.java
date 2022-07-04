package com.emailProject.emailProject.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonManagedReference;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.*;
import java.util.ArrayList;
import java.util.List;

@Entity
@Data
@Table(name = "course")
@NoArgsConstructor
public class Course {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "course_id")
    private Long id;

    @Column(name = "name")
    private String className;

    @Column(name = "teacher")
    private String teacherName;

    @OneToMany(mappedBy = "course")
    private List<Student> students;

    public Course(String className, String teacherName){
        this.className = className;
        this.teacherName = teacherName;
    }

    public Course(Long id, String className, String teacherName, List<Student> students) {
        this.id = id;
        this.className = className;
        this.teacherName = teacherName;
        this.students = students;
    }
}
