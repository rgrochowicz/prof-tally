import requests
from bs4 import BeautifulSoup as bs
from bs4 import NavigableString
import re
import psycopg2
import time

class Course(object):

	def __init__(self, node):

		#get all children that aren't strings
		children = [x for x in node.children if not isinstance(x, NavigableString)]

		self.crn = int(children[0].a.string)
		self.max = int(children[12].string)
		self.enrolled = int(children[15].string)

def get_stats():

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

	r = requests.post("http://banner.rowan.edu/reports/reports.pl?task=Section_Tally", data=datas)
	text = r.text

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

		#put the course stat into the db
		cur.execute("INSERT INTO course_enrollments (crn, max, enrolled, time) VALUES (%s, %s, %s, date_trunc('minute', now()))",
			(course.crn, course.max, course.enrolled))

	#commit and close connection
	conn.commit()
	cur.close()
	conn.close()

	print("Updated at: {}".format(time.strftime('%X %x %Z')))


def main():

	#loop forever getting stats and waiting an hour
	while True:
		get_stats()
		time.sleep(60 * 60)



if __name__ == '__main__':
	main()