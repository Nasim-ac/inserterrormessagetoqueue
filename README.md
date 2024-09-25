# inserterrormessagetoqueue

This package takes a db connection(likely CDR),subject,body to insert a error message to dbo.MessageQueue only if the same subject+body doesn't exist on the same day. The code can be reused in any service which runs in an interval and retry agian and again if  execution fails in each repetation.