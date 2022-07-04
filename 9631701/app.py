# flask imports
from flask import Flask, request, jsonify, make_response
from flask_sqlalchemy import SQLAlchemy
import uuid # for public id
from werkzeug.security import generate_password_hash, check_password_hash
# imports for PyJWT authentication
import jwt
from datetime import datetime, timedelta
from functools import wraps

import re
from flask_mail import Mail, Message
import threading

# creates Flask object
app = Flask(__name__)
# configuration
# NEVER HARDCODE YOUR CONFIGURATION IN YOUR CODE
# INSTEAD CREATE A .env FILE AND STORE IN IT
app.config['SECRET_KEY'] = '004f2af45d3a4e161a7dd2d17fdae47f'
# database name
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///Database.db'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = True
# creates SQLALCHEMY object
db = SQLAlchemy(app)

app.config['MAIL_SERVER'] = 'smtp.gmail.com'
app.config['MAIL_PORT'] = 465
app.config['MAIL_USE_TLS'] = False
app.config['MAIL_USE_SSL'] = True
app.config['MAIL_USERNAME'] = 'rasoul.kh41@gmail.com'
app.config['MAIL_PASSWORD'] = 'vbhsrsoqnheigcwj'

mail = Mail(app)

@app.before_first_request
def create_tables():
    db.create_all()
    app.preprocess_request()

class ClassRoomORM(db.Model):
    __tablename__ = 'classroom'
    id = db.Column(db.Integer, primary_key=True)
    classID = db.Column(db.String(20), unique=True)
    Professor = db.Column(db.String(20))

class StudentModel(db.Model):
    __tablename__ = 'students'
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(80))
    # check mail @domain.com for domain name (gmail, yahoo, etc)
    email = db.Column(db.String(80), unique=True)
    classID = db.Column(db.String(80))
    # score in range of 0-20 
    score = db.Column(db.Integer)

class User(db.Model):
    id = db.Column(db.Integer, primary_key = True)
    public_id = db.Column(db.String(50), unique = True)
    name = db.Column(db.String(100))
    email = db.Column(db.String(70), unique = True)
    password = db.Column(db.String(80))

# decorator for verifying the JWT
def token_required(f):
    @wraps(f)
    def decorator(*args, **kwargs):
        token = None
        if 'x-access-tokens' in request.headers:
            token = request.headers['x-access-tokens']
            print(token)
        if not token:
            return jsonify({'message': 'a valid token is missing'})
        try:
            data = jwt.decode(token, app.config['SECRET_KEY'], algorithms=["HS256"])
            print(data)
            current_user = User.query.filter_by(public_id=data['public_id']).first()
        except:
            return jsonify({'message': 'token is invalid'})
 
        return f(current_user, *args, **kwargs)
    return decorator

# User Database Route
# this route sends back list of users
@app.route('/user', methods =['GET'])
@token_required
def get_all_users(current_user):
	# querying the database
	# for all the entries in it
	users = User.query.all()
	# converting the query objects
	# to list of jsons
	output = []
	for user in users:
		# appending the user data json
		# to the response list
		output.append({
			'public_id': user.public_id,
			'name' : user.name,
			'email' : user.email
		})

	return jsonify({'users': output})

# route for logging user in
@app.route('/login', methods =['POST'])
def login():
	# creates dictionary of form data
	auth = request.form

	if not auth or not auth.get('email') or not auth.get('password'):
		# returns 401 if any email or / and password is missing
		return make_response(
			'Could not verify',
			401,
			{'WWW-Authenticate' : 'Basic realm ="Login required !!"'}
		)

	user = User.query\
		.filter_by(email = auth.get('email'))\
		.first()

	if not user:
		# returns 401 if user does not exist
		return make_response(
			'Could not verify',
			401,
			{'WWW-Authenticate' : 'Basic realm ="User does not exist !!"'}
		)

	if check_password_hash(user.password, auth.get('password')):
		# generates the JWT Token
		token = jwt.encode({
			'public_id': user.public_id,
			'exp' : datetime.utcnow() + timedelta(minutes = 30)
		}, app.config['SECRET_KEY'])
        

		return make_response(jsonify({'token' : token}), 201)
	# returns 403 if password is wrong
	return make_response(
		'Could not verify',
		403,
		{'WWW-Authenticate' : 'Basic realm ="Wrong Password !!"'}
	)

# signup route
@app.route('/signup', methods =['POST'])
def signup():
	# creates a dictionary of the form data
	data = request.form

	# gets name, email and password
	name, email = data.get('name'), data.get('email')
	password = data.get('password')

	# checking for existing user
	user = User.query\
		.filter_by(email = email)\
		.first()
	if not user:
		# database ORM object
		user = User(
			public_id = str(uuid.uuid4()),
			name = name,
			email = email,
			password = generate_password_hash(password)
		)
		# insert user
		db.session.add(user)
		db.session.commit()

		return make_response('Successfully registered.', 201)
	else:
		# returns 202 if user already exists
		return make_response('User already exists. Please Log in.', 202)


@app.route('/student/<string:email>', methods=['GET'])
@token_required
def get_student(current_user, email):
    student = StudentModel.query.filter_by(email=email).first()

    if not student:
        return jsonify({'message': 'No student found!'})

    return jsonify({'student': student.name})


@app.route('/student/<string:email>', methods=['PUT'])
@token_required
def update_student(current_user, email):
    student = StudentModel.query.filter_by(email=email).first()

    if not student:
        return jsonify({'message': 'No student found!'})

    student.name = request.json['name']
    student.email = request.json['email']
    student.classID = request.json['classID']
    student.score = request.json['score']

    db.session.commit()

    return jsonify({'message': 'Student updated!'})


@app.route('/student', methods=['POST'])
@token_required
def create_student(current_user):
    try:
        data = request.form
        # input can be a list of students
        if not isinstance(data, list):
            data = [data]
        print(data)
        for student in data:
            
            # check email is valid with regex
            if not re.match(r"[^@]+@[^@]+\.[^@]+", student['email']):
                return jsonify({'message': 'Invalid email address!'})
            # check score is integer and between 0 and 20
            if not student['score'].isdigit() or int(student['score']) < 0 or int(student['score']) > 20:
                return jsonify({'message': 'Invalid score!'})
            new_student = StudentModel(name=request.form.get('name'), email=request.form.get('email'), classID=request.form.get('classID'), score=int(request.form.get('score')))
            db.session.add(new_student)
            db.session.commit()
        


        return jsonify({'message': 'Student created!'})
    except:
        return jsonify({'message': 'Something went wrong!'})


@app.route('/student/<string:email>', methods=['DELETE'])
@token_required
def delete_student(current_user, email):
    student = StudentModel.query.filter_by(email=email).first()

    if not student:
        return jsonify({'message': 'No student found!'})

    db.session.delete(student)
    db.session.commit()

    return jsonify({'message': 'Student deleted!'})


@app.route('/class/<string:classID>', methods=['GET'])
@token_required
def get_class(current_user, classID):
    class_ = ClassRoomORM.query.filter_by(classID=classID).first()

    if not class_:
        return jsonify({'message': 'No class found!'})

    return jsonify({'class': class_.Professor})

# create class
@app.route('/class', methods=['POST'])
@token_required
def create_class(current_user):
    data = request.form
        
    # check class name is not empty
    new_class = ClassRoomORM(classID=request.form.get('classID'),
                            Professor=request.form.get('professor'))
    db.session.add(new_class)
    db.session.commit()
    


    return jsonify({'message': 'Class created!'})



@app.route('/class/<string:classID>', methods=['GET'])
@token_required
def get_class_student(current_user, classID):
    student = StudentModel.query.filter_by(classID=classID).all()

    if not student:
        return jsonify({'message': 'No student found!'})

    return jsonify({'student': [student.name for student in student]})

@app.route('/class/<string:classID>', methods=['PUT'])
@token_required
def update_class(current_user, classID):
    class_ = ClassRoomORM.query.filter_by(classID=classID).first()

    if not class_:
        return jsonify({'message': 'No class found!'})

    class_.classID = request.json['classID']
    class_.Professor = request.json['professor']


    db.session.commit()

    return jsonify({'message': 'Class updated!'})


@app.route('/class/<string:classID>', methods=['DELETE'])
@token_required
def delete_class(current_user, classID):
    class_ = ClassRoomORM.query.filter_by(classID=classID).first()

    if not class_:
        return jsonify({'message': 'No class found!'})

    db.session.delete(class_)
    db.session.commit()

    return jsonify({'message': 'Class deleted!'})


# Rest API for send mail to classID students
@app.route('/send_mail/<string:classID>', methods=['POST'])
@token_required
def send_mail(current_user, classID):
    def send_email(name, email, classID, score):
        with app.test_request_context():
            from flask import request
            msg = Message(
                'Hello ',
                sender='Rasoul Khazaei',
                recipients=[email]
            )
            msg.body = 'Hello ' + name + '!\n' +  'message' +'\n'+'Your grade is '+str(score) + 'in the class ' + classID + 'with the professor ' + ClassRoomORM.query.filter_by(classID=classID).first().Professor
            mail.send(msg)
            return jsonify({'message': 'Email for ' + name + ' sent!'})

    student = StudentModel.query.filter_by(classID=classID).all()
    
    if not student:
        return jsonify({'message': 'No student found!'})
    return_value = []
    for s in student:
        # send email with threading
        with app.app_context():
            # send_email(s.name, s.email, classID, s.score)
            th = threading.Thread(target=send_email, args=(s.name, s.email, classID, s.score))
            th.start()
            return_value.append(s.name)
    
    return jsonify({'message': 'Email for ' + str(return_value) + ' sent!'})





if __name__ == '__main__':
    app.run(debug=True)
    # app.run(host='