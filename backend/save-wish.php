<?php
   header('Content-Type: application/json');
   header('Access-Control-Allow-Origin: https://dark-logos.github.io');
   header('Access-Control-Allow-Methods: POST');
   header('Access-Control-Allow-Headers: Content-Type');

   // Получаем данные из POST-запроса
   $input = file_get_contents('php://input');
   $data = json_decode($input, true);

   if (!isset($data['wish']) || empty($data['wish'])) {
       echo json_encode(['success' => false, 'error' => 'Желание не указано']);
       exit;
   }

   // Санитизация ввода
   $wish = htmlspecialchars($data['wish'], ENT_QUOTES, 'UTF-8');

   // Запись в текстовый файл
   $file = 'wishes.txt';
   $timestamp = date('Y-m-d H:i:s');
   $entry = "[$timestamp] $wish\n";

   if (file_put_contents($file, $entry, FILE_APPEND | LOCK_EX) !== false) {
       echo json_encode(['success' => true]);
   } else {
       echo json_encode(['success' => false, 'error' => 'Не удалось сохранить желание']);
   }
   ?>