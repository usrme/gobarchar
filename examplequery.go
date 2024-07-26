package gobarchar

import (
	"fmt"
	"strings"
)

type query struct {
	Title  string
	Keys   []string
	Values []string
}

type exampleQueries []query

// Using 'map[string]string' would be more readable, but would require additional logic
// to keep insertion order intact
var (
	presidentQueryTitle string = "Presidents of the United States by age at start of presidency"
	presidentQuery             = query{
		Title: presidentQueryTitle,
		Keys: []string{
			"George Washington (1789)", "John Adams (1797)", "Thomas Jefferson (1801)", "James Madison (1809)", "James Monroe (1817)", "John Quincy Adams (1825)", "Andrew Jackson (1829)", "Martin Van Buren (1837)", "William Henry Harrison (1841)", "John Tyler (1841)", "James K. Polk (1845)", "Zachary Taylor (1849)", "Millard Fillmore (1850)", "Franklin Pierce (1853)", "James Buchanan (1857)", "Abraham Lincoln (1861)", "Andrew Johnson (1865)", "Ulysses S. Grant (1869)", "Rutherford B. Hayes (1877)", "James A. Garfield (1881)", "Chester A. Arthur (1881)", "Grover Cleveland (first term) (1885)", "Benjamin Harrison (1889)", "Grover Cleveland (second term) (1893)", "William McKinley (1897)", "Theodore Roosevelt (1901)", "William Howard Taft (1909)", "Woodrow Wilson (1913)", "Warren G. Harding (1921)", "Calvin Coolidge (1923)", "Herbert Hoover (1929)", "Franklin D. Roosevelt (1933)", "Harry S. Truman (1945)", "Dwight D. Eisenhower (1953)", "John F. Kennedy (1961)", "Lyndon B. Johnson (1963)", "Richard Nixon (1969)", "Gerald Ford (1974)", "Jimmy Carter (1977)", "Ronald Reagan (1981)", "George H. W. Bush (1989)", "Bill Clinton (1993)", "George W. Bush (2001)", "Barack Obama (2009)", "Donald Trump (2017)", "Joe Biden (2021)", "spaces", "sort", "title",
		},
		Values: []string{
			"57", "61", "57", "57", "58", "57", "61", "54", "68", "51", "49", "64", "50", "48", "65", "52", "56", "46", "54", "49", "51", "47", "55", "55", "54", "42", "51", "56", "55", "51", "54", "51", "60", "62", "43", "55", "56", "61", "52", "69", "64", "46", "54", "47", "70", "78", "yes", "asc", presidentQueryTitle,
		},
	}
	terraformQueryTitle string = "% of community PRs opened against Terraform after license change (2023)"
	terraformQuery             = query{
		Title: terraformQueryTitle,
		Keys: []string{
			"February", "March", "April", "May", "June", "July", "August", "September", "title",
		},
		Values: []string{
			"14.71", "22.83", "28.57", "22.58", "23.29", "20.29", "9.30", "9.52", terraformQueryTitle,
		},
	}
	Examples = exampleQueries{presidentQuery, terraformQuery}
)

func CreateListItems(url string, elements exampleQueries) string {
	start := `<hr/>
	<p>More example queries:</p>
	<ul>
`
	end := `
	</ul>
`
	var listItems strings.Builder
	for _, v := range elements {
		items := make([]string, 0, len(v.Keys))
		for i, j := range v.Keys {
			items = append(items, fmt.Sprintf("%s=%s", j, v.Values[i]))
		}
		listItems.WriteString(
			fmt.Sprintf("<li><a href='%s?%s'>%s</a></li>", url, strings.Join(items, "&"), v.Title),
		)
		items = nil
	}
	return fmt.Sprintf("%s%s%s", start, listItems.String(), end)
}
