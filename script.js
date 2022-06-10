import http from 'k6/http';

export default function () {
  http.options('http://localhost:26726');
}
