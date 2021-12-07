# Article

The definition of all fields from an article.

* title : String!
  *  of course the title
* subtitle : String!
  * the subtitle shown under the title
* content_md : String!
  * The content field which is filled in the admin page and will converted while insert or update the article into the database to html
* content : String!
  * the generated content from the content_md field. It is generated while inserting or updating
* created_date : String!
  * the date from the creation of the article. It is in following german format: DD.MM.YYYY HH:mm
* modified_date : String
  * the date from the last modified date of the article. Empty String if it is not modified. It is in following german format: DD.MM.YYYY HH:mm
* tags : [String]!
  * the list of tags from the article
* categories : [String]!
  * the categories from the article. The difference from tags is that categories are more generic and so specific than tags
* author : [String]!
  * the author from the article
* needsTOC : Boolean!
  * A bool that shows if the article needs a table of content which is generated in the backend before inserted it in the database
* toc : String
  * the generated from toc 
* isInSeries: Boolean!
  * if the article is in a series (a series or a specific order from articles that you should read to follow a topic) then you can speficy here that it is in it. 
* series : String
  * the series name


## legend
NAME_OF_FIELD : TYPE (! means the field is not optional and never null or undefined)