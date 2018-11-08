package merger_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/johandry/merger"
)

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type Grade struct {
	Teacher string  `json:"teacher" mapstructure:"teacher"`
	Number  float32 `json:"number" mapstructure:"number"`
}

type Student struct {
	Name      string           `json:"name" mapstructure:"name"`
	TextBooks []string         `json:"text_books" mapstructure:"text_books"`
	Address   Address          `json:"address" mapstructure:"address"`
	Grades    map[string]Grade `json:"grades" mapstructure:"grades"`
	GPA       float32          `json:"gpa" mapstructure:"gpa"`
}

// RefreshGPA recalculate the student GPA from the existing grades
func (s *Student) RefreshGPA() {
	var total float32
	for _, grade := range s.Grades {
		total = total + grade.Number
	}
	s.GPA = total / float32(len(s.Grades))
}

var configFromFile string

func init() {
	// Do this in shell using `export` command. This is to mock environment variables
	// i.e. export EXAMPLE_name=John
	os.Setenv("EXAMPLE_name", "John")
	os.Setenv("EXAMPLE_books", "B1,B2, B3, 'The Book', B4")
	os.Setenv("EXAMPLE_address__City", "San Diego")
	os.Setenv("EXAMPLE_address__Country", "US")

	configFromFile = `{"name": "Mary", "address": {"city": "San Diego", "country": "US"}}`
}

func Example() {
	// Load information from environment variables
	studentEnvVars := make(map[string]string, 0)
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "EXAMPLE_") {
			continue
		}
		envVar := strings.Split(env, "=")
		name := strings.TrimLeft(envVar[0], "EXAMPLE_")
		studentEnvVars[name] = envVar[1]
	}

	// Load information from configuration file
	studentFromFile := Student{}
	json.Unmarshal([]byte(configFromFile), &studentFromFile)

	// Information from code
	studentWithGrades := Student{
		Grades: map[string]Grade{
			"Science": Grade{
				Teacher: "Dr. Smith",
				Number:  89.99,
			},
			"Computer Science": Grade{
				Teacher: "Dr. Steve",
				Number:  99.99,
			},
		},
	}

	student := Student{}

	if err := merger.Merge(&student, studentEnvVars, studentFromFile, studentWithGrades); err != nil {
		log.Fatal(fmt.Errorf("Failed to merge the studens information. %s", err))
	}

	student.RefreshGPA()

	fmt.Printf("Name: %s, Books: %s, Address: %+v, Computer Science Grade: %+v, Science Grade: %+v, GPA: %f",
		student.Name,
		student.TextBooks,
		student.Address,
		student.Grades["Computer Science"],
		student.Grades["Science"],
		student.GPA,
	)
	// Output: Name: John, Books: [], Address: {City:San Diego Country:US}, Computer Science Grade: {Teacher:Dr. Steve Number:99.99}, Science Grade: {Teacher:Dr. Smith Number:89.99}, GPA: 94.989998
}
