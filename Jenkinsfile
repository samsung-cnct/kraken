podTemplate(label: 'k2cli', containers: [
   containerTemplate(name: 'golang', image: 'golang:1.7.5', ttyEnabled: true, command: 'cat')
 ]) {
   node('k2cli') {
       container('golang'){

           stage('hello!') {
               echo 'hello world!'
           }
         }
       }
     }
