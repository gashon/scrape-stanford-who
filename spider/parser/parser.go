package parser

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Person struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Affiliation string `json:"affiliation"`
	Department string `json:"department"`
	Role string `json:"role"`
	ProfileURL string `json:"profile_url"`
}

type ElementNodeLayout struct {
	Data string
	Class   string
	NextNode   *ElementNodeLayout
}

func containsClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && attr.Val == class {
			return true
		}
	}
	return false
}

func recursivelyGetTargetNode(parent *html.Node, layout *ElementNodeLayout) *html.Node {
	if layout == nil {
		return parent
	}

	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == layout.Data && (len(layout.Class) == 0 || containsClass(c, fmt.Sprintf(`\"%s\"`, layout.Class))) {
			return recursivelyGetTargetNode(c, layout.NextNode)
		}
	}

	return nil
}

func findChildNodeByLayout(parent *html.Node, layout *ElementNodeLayout) *html.Node {
	return recursivelyGetTargetNode(parent, layout)
}

func findChildNodeByLayouts(parent *html.Node, layouts []*ElementNodeLayout) *html.Node {
	for _, layout := range layouts {
		if node := findChildNodeByLayout(parent, layout); node != nil {
			return node
		}
	}
	return nil
}

func Parse(htmlContent string) []Person {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}

	people := []Person{}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" && containsClass(n, `\"people__list__item\"`) {
			personalInfo := make(map[string]string)
			
			personalAttributesNode := findChildNodeByLayout(n, &ElementNodeLayout{
				Data: "div",
				Class: "person__data",
				NextNode: &ElementNodeLayout{
					Data: "div",
					Class: "person__data__main",
				},
			});

			// name
			nameNode := findChildNodeByLayout(personalAttributesNode, &ElementNodeLayout{
				Data: "a",
				Class: "person__data__name",
				NextNode: &ElementNodeLayout{
					Data: "span",
					Class: "",
				},
			});

			if nameNode != nil {
				personalInfo["name"] = nameNode.FirstChild.Data
			}

			// affiliation
			affiliationNode := findChildNodeByLayout(personalAttributesNode, &ElementNodeLayout{
				Data: "div",
				Class: "person__data__affiliation",
				NextNode: &ElementNodeLayout{
					Data: "span",
					Class: "",
				},
			});

			if affiliationNode != nil {
				personalInfo["affiliation"] = affiliationNode.FirstChild.Data
			}

			// department
			departmentNode := findChildNodeByLayout(personalAttributesNode, &ElementNodeLayout{
				Data: "div",
				Class: "person__data__department",
				NextNode: &ElementNodeLayout{
					Data: "span",
					Class: "",
				},
			});

			if departmentNode != nil {
				personalInfo["department"] = departmentNode.FirstChild.Data
			}

			// role
			roleNode := findChildNodeByLayout(personalAttributesNode, &ElementNodeLayout{
				Data: "div",
				Class: "person__data__role",
				NextNode: &ElementNodeLayout{
					Data: "span",
					Class: "",
				},
			});

			if roleNode != nil {
				personalInfo["role"] = roleNode.FirstChild.Data
			}

			// profileURL
			profileURLNode := findChildNodeByLayout(n, &ElementNodeLayout{
				Data: "div",
				Class: "person__data",
				NextNode: &ElementNodeLayout{
					Data: "div",
					Class: "person__data__profile__pic__container",
					NextNode: &ElementNodeLayout{
						Data: "img",
						Class: "",
					},
				},
			});

			if profileURLNode != nil {
				for _, attr := range profileURLNode.Attr {
					if attr.Key == "src" {
						personalInfo["profileURL"] = attr.Val
						break;
					}
				}
			}
			
			// email
			emailNode := findChildNodeByLayout(n, &ElementNodeLayout{
				Data: "ul",
				Class: "",
				NextNode: &ElementNodeLayout{
					Data: "li",
					Class: "person__vcard__list__item",
					NextNode: &ElementNodeLayout{
						Data: "span",
						Class: "sidebar-detail",
					},
				},
			});

			if emailNode != nil {
				for _, attr := range emailNode.FirstChild.FirstChild.Attr {
					if attr.Key == "title" {
						personalInfo["email"] = attr.Val
						break;
					}
				}
			}

			people = append(people, Person{
				Name: personalInfo["name"],
				Affiliation: personalInfo["affiliation"],
				Department: personalInfo["department"],
				Role: personalInfo["role"],
				ProfileURL: personalInfo["profileURL"],
				Email: personalInfo["email"],
			})

		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}

	}

	traverse(doc)

	return people
}