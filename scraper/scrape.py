import requests
from bs4 import BeautifulSoup as bs
from bs4 import NavigableString
import re
import psycopg2
import json
import os.path

class Course(object):

	def __init__(self, node):

		#get all children that aren't strings
		children = [x for x in node.children if not isinstance(x, NavigableString)]

		self.crn = int(children[0].a.string)
		self.subject = str(children[1].string)
		self.course_num = str(children[2].string)
		self.section = str(children[3].string)
		# 4 is "part of term"
		# 5 is "session" (time of day)
		self.title = str(list(children[6].children)[0].string)
		self.professors = list(map(str, children[7].stripped_strings))
		self.raw_times = list(map(str, children[8].stripped_strings))
		self.campus = str(children[9].string)
		self.hours = float(str(children[11].string).strip())
		self.max = int(children[12].string)
		self.max_reserved = int(children[13].string)
		self.left_reserved = int(children[14].string)
		self.enrolled = int(children[15].string)
		self.available = int(children[16].string)

		self.times = []
		for time in self.raw_times:
			self.times.append(CourseTime(time))

class CourseTime(object):

	#ugh. regexes...
	matcher = re.compile(r"(?P<days>[MTWRFS]+)(?:\s*)(?P<start_time>\d{4})?(?:\s*)(?P<end_time>\d{4})(?:\s*)(?P<building>.*?)(?:\s+)(?P<room>.*?)(?:\s*\()(?P<type>.+?)(?:\))")

	def __init__(self, raw_time):
		match = CourseTime.matcher.match(raw_time)
		self.raw_time = raw_time
		self.invalid = True

		if match:
			self.days = match.group('days')
			self.start_time = match.group('start_time')
			self.end_time = match.group('end_time')
			self.building = match.group('building')
			self.room = match.group('room')
			self.type = match.group('type')
			self.invalid = False

class Attribute(object):

	def __init__(self, node):
		self.short = node['value']

		#the names come out as "short - name", so the names are extracted using regex
		match = re.match(r'.*?\s-\s(.*)', str(node.string))
		if match:
			self.name = match.group(1)

	def __str__(self):
		return "{} - {}".format(self.short, self.name)

def get_attributes():

	#connect to postgres
	conn = psycopg2.connect(database=os.environ['POSTGRES_DATABASE'],
		user=os.environ['POSTGRES_USER'],
		password=os.environ['POSTGRES_PASSWORD'],
		host=os.environ['POSTGRES_HOST'],
		port=os.environ['POSTGRES_PORT'])
	cur = conn.cursor()

	payload = {
		"Select Term": "Select Term",
		"term": "201520",
		"task": "Section_Tally"
	}
	r = requests.post("http://banner.rowan.edu/reports/reports.pl?task=Section_Tally", data=payload)
	soup = bs(r.text)

	#this uses the attribute dropdown to get its values
	attribute_select = soup.select("select[name=attr]")[0]
	course_attributes = []

	for option in attribute_select.select("option"):
		attr = Attribute(option)

		#ignore the ALL attribute
		if attr.short != 'ALL':
			course_attributes.append(attr)

	for course_attribute in course_attributes:
		cur.execute(
			"""INSERT INTO course_attrs (
					short,
					name
				) VALUES (%s, %s)
			""", (
				course_attribute.short,
				course_attribute.name
			))
	conn.commit()


	#loop through and do a separate request for each attribute
	for course_attribute in course_attributes:
		attr_search = {
			"term": "201520",
			"task": "Section_Tally",
			"coll": "ALL",
			"dept": "ALL",
			"subj": "ALL",
			"ptrm": "ALL",
			"sess": "ALL",
			"prof": "ALL",
			"attr": course_attribute.short,
			"camp": "ALL",
			"bldg": "ALL",
			"Search": "Search"
		}

		for course in get_courses_in_page(requests.post("http://banner.rowan.edu/reports/reports.pl?task=Section_Tally", data=attr_search).text):

			cur.execute(
				"""INSERT INTO courses_and_attrs (
						crn,
						attr
					) VALUES (%s, %s)
				""", (
					course.crn,
					course_attribute.short
				))

	#commit and close connection
	conn.commit()
	cur.close()
	conn.close()

#given page html, return the courses
def get_courses_in_page(text):

	text = text.replace('></div><div>', '>') #fix html error

	soup = bs(text) 
	class_table = soup.find("table", class_="report")

	courses = []

	#loop through all table rows
	for tr in class_table.tbody.children:

		#skip strings and a tags
		if isinstance(tr, NavigableString): continue
		if tr.name == "a": continue

		#make sure it's the correct type of row
		if 'style' in tr.attrs and tr.attrs['style'] == 'background-color:inherit':

			#also make sure it's not cancelled
			if 'class' not in tr.attrs or tr.attrs['class'] != 'cancel':
				course = Course(tr)
				courses.append(course)

	return courses

def main():

	#get everything from the section tally
	datas = {
		"term": "201520",
		"task": "Section_Tally",
		"coll": "ALL",
		"dept": "ALL",
		"subj": "ALL",
		"ptrm": "ALL",
		"sess": "ALL",
		"prof": "ALL",
		"attr": "ALL",
		"camp": "ALL",
		"bldg": "ALL",
		"Search": "Search"
	}

	# download and cache in page.html
	text = ""
	if os.path.isfile('page.html'):
		with open('page.html') as page:
			text = page.read()
	else:
		r = requests.post("http://banner.rowan.edu/reports/reports.pl?task=Section_Tally", data=datas)
		text = r.text
		with open('page.html', 'w') as page:
			page.write(text)


	courses = get_courses_in_page(text)

	conn = psycopg2.connect(database=os.environ['POSTGRES_DATABASE'],
		user=os.environ['POSTGRES_USER'],
		password=os.environ['POSTGRES_PASSWORD'],
		host=os.environ['POSTGRES_HOST'],
		port=os.environ['POSTGRES_PORT'])

	cur = conn.cursor()
	for course in courses:

		#put the course into the db
		cur.execute(
			"""INSERT INTO courses (
				crn,
				subject,
				course_num,
				section,
				title,
				professors,
				campus,
				hours,
				max,
				max_reserved,
				left_reserved,
				enrolled,
				available
			)
			VALUES (%s, %s, %s, %s, %s, %s,
					%s, %s, %s, %s, %s, %s, %s)""",
			(	course.crn,
				course.subject,
				course.course_num,
				course.section,
				course.title,
				course.professors,
				course.campus,
				course.hours,
				course.max,
				course.max_reserved,
				course.left_reserved,
				course.enrolled,
				course.available	)
		)

		for time in course.times:
			if time.invalid:
				#insert a row to fix later
				cur.execute(
					"""INSERT INTO course_times (
						course_crn,
						weekday,
						start_time,
						length,
						building,
						room,
						type,
						invalid,
						raw_time
					)
					VALUES (%(course_crn)s, null, null, null, null,
							null, null, %(invalid)s, %(raw_time)s)""",
					{
						"course_crn": course.crn,
						"invalid": time.invalid,
						"raw_time": time.raw_time
					}
				)
			else:
				#insert a row per day
				for day in time.days:
					cur.execute(
						"""INSERT INTO course_times (
							course_crn,
							weekday,
							start_time,
							length,
							building,
							room,
							type,
							invalid,
							raw_time
						)
						VALUES (%(course_crn)s, %(day)s, %(start_time)s,
								%(end_time)s::time - %(start_time)s::time,
								%(building)s, %(room)s,  %(type)s,
								%(invalid)s, %(raw_time)s)""",
						{
							"course_crn": course.crn,
							"day": day,
							"start_time": time.start_time,
							"end_time": time.end_time,
							"building": time.building,
							"room": time.room,
							"type": time.type,
							"invalid": time.invalid,
							"raw_time": time.raw_time
						}
					)

	#commit and close connection
	conn.commit()
	cur.close()
	conn.close()

	#get and associate attributes with courses
	get_attributes()

if __name__ == '__main__':
	main()

# vim:tabstop=4:noexpandtab
