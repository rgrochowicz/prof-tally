import requests
from bs4 import BeautifulSoup as bs
from bs4 import NavigableString
import re
import psycopg2
import json

class Course(object):

	def __init__(self, node):

		#get all children that aren't strings
		children = [x for x in node.children if not isinstance(x, NavigableString)]

		self.crn = int(children[0].a.string)
		self.subject = str(children[1].string)
		self.course_num = str(children[2].string)
		self.section = str(children[3].string)
		# 4 is "part of term", 5 is "session" (time of day)
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


def main():


	#set to true to download fresh from rowan
	download_fresh = False

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

	text = ""
	if download_fresh:
		r = requests.post("http://banner.rowan.edu/reports/reports.pl?task=Section_Tally", data=datas)
		text = r.text
	else:
		with open('page.html') as page:
			text = page.read()


	courses = []

	soup = bs(text)
	class_table = soup.find("table", class_="report")

	#loop through all table rows
	for tr in class_table.tbody.children:
		#skip strings and a tags
		if isinstance(tr, NavigableString): continue
		if tr.name == "a": continue

		#make sure it's the correct type of row
		if 'style' in tr.attrs and tr.attrs['style'] == 'background-color:inherit':
			courses.append(Course(tr))

	conn = psycopg2.connect("dbname=tally user=tally password=tally host=127.0.0.1")
	cur = conn.cursor()
	for course in courses:

		#put the course into the db
		cur.execute(
			"INSERT INTO courses (				\
				crn,							\
				subject,						\
				course_num,						\
				section,						\
				title,							\
				professors,						\
				campus,							\
				hours,							\
				max,							\
				max_reserved,					\
				left_reserved,					\
				enrolled,						\
				available						\
			)									\
			VALUES (%s, %s, %s, %s, %s, %s,		\
					%s, %s, %s, %s, %s, %s, %s) \
			RETURNING id",
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
		course_id = cur.fetchone()[0]

		for time in course.times:
			if time.invalid:
				#insert a row to fix later
				cur.execute(
					"INSERT INTO course_times (						\
						course_id,									\
						weekday,									\
						start_time,									\
						length,										\
						building,									\
						room,										\
						type,										\
						invalid,									\
						raw_time									\
					)												\
					VALUES (%(course_id)s, null, null, null, null,	\
							null, null, %(invalid)s, %(raw_time)s)",
					{	"course_id": course_id,
						"invalid": time.invalid,
						"raw_time": time.raw_time	}
				)
			else:
				#insert a row per day
				for day in time.days:
					cur.execute(
						"INSERT INTO course_times (		\
							course_id,					\
							weekday,					\
							start_time,					\
							length,						\
							building,					\
							room,						\
							type,						\
							invalid,					\
							raw_time					\
						)								\
						VALUES (%(course_id)s, %(day)s, %(start_time)s,			\
								%(end_time)s::time - %(start_time)s::time,		\
								%(building)s, %(room)s,  %(type)s,				\
								%(invalid)s, %(raw_time)s)",
						{	"course_id": course_id,
							"day": day,
							"start_time": time.start_time,
							"end_time": time.end_time,
							"building": time.building,
							"room": time.room,
							"type": time.type,
							"invalid": time.invalid,
							"raw_time": time.raw_time	}
					)

	#commit and close connection
	conn.commit()
	cur.close()
	conn.close()

if __name__ == '__main__':
	main()

# vim:tabstop=4:noexpandtab
