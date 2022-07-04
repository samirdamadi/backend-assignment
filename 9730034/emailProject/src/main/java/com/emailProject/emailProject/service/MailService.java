package com.emailProject.emailProject.service;

import com.emailProject.emailProject.model.Course;
import com.emailProject.emailProject.model.Student;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.Map;
@Service
public class MailService {
    @Autowired
    private JavaMailSender javaMailSender;

    @Value("${spring.mail.username}")
    private String sender;

    public Map<Student, Boolean> sendMessages(Course course){
        Map<Student, Boolean> mailResults = new HashMap<>();
        for (Student s:course.getStudents()
             ) {

            SendMailThread sendMailThread = new SendMailThread(course.getClassName(), course.getTeacherName(), s);
            Thread thread = new Thread(sendMailThread);
            thread.start();
            try {
                thread.join();
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            mailResults.put(s, sendMailThread.getMailResult());
        }
        return mailResults;
    }


    class SendMailThread implements Runnable{
        private String className;
        private String teacherName;
        private Student student;
        private Boolean mailResult;
        public SendMailThread(String className, String teacherName, Student student){
            this.className = className;
            this.student = student;
            this.teacherName = teacherName;
        }
        @Override
        public void run() {
            try {
                SimpleMailMessage mailMessage = new SimpleMailMessage();
                mailMessage.setFrom(sender);
                mailMessage.setTo(student.getEmail());
                mailMessage.setText("className : " + className +
                        "\nteacherName : " + teacherName +
                        "\nstudentName :" + student.getName() +
                        "\nstudent's score : " + student.getScore());
                javaMailSender.send(mailMessage);
                mailResult =true;
            }catch (Exception e){
                mailResult = false;
            }
        }
        public Boolean getMailResult(){
            return mailResult;
        }
    }
}
