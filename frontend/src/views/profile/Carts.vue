<script setup>

</script>

<template>
  <div class="content-header">
    <div class="container-fluid">
      <div class="row mb-2">
        <div class="col-sm-6">
          <h1 class="m-0">
            Товары                    </h1>
        </div><!-- /.col -->
        <div class="col-sm-6">
        </div><!-- /.col -->
      </div><!-- /.row -->
    </div><!-- /.container-fluid -->
  </div>
  <div class="content">

<!--    <script src="https://code.jquery.com/jquery-3.7.1.js" integrity="sha256-eKhayi8LEQwp4NKxN+CfCh+3qOVUtJn3QNZ0TciWLP4=" crossorigin="anonymous"></script>-->

    <div>
      <!--    <h2 style="font-size: 18px;">Запросить у WB обновление: <a href="--><!--"-->
      <!--                                                               class="bi bi-arrow-clockwise bold"></a></h2>-->
      <h2 style="font-size: 18px;">Запросить у WB обновление: <span class="bi bi-arrow-clockwise bold wbarticles " style="color: green; cursor: pointer; transform-origin: center;" ;=""></span></h2>
      <!--    -->    <div id="w0" class="grid-view"><div class="summary">Показаны записи <b>1-9</b> из <b>9</b>.</div>
      <table class="table table-striped table-bordered"><thead>
      <tr><th>#</th><th>Фото</th><th><a href="/artlist?sort=articule" data-sort="articule">Артикул</a></th><th><a class="asc" href="/artlist?sort=-name" data-sort="-name">Название</a></th><th><a href="/artlist?sort=created" data-sort="created">Создан</a></th><th><a href="/artlist?sort=updated" data-sort="updated">Последнее обновление</a></th><th><a href="/artlist?sort=cost_price" data-sort="cost_price">Себестоимость</a></th><th>Действие</th></tr><tr id="w0-filters" class="filters"><td>&nbsp;</td><td>&nbsp;</td><td><input type="text" class="form-control" name="WbArticles[articule]"></td><td><input type="text" class="form-control" name="WbArticles[name]"></td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td></tr>
      </thead>
        <tbody>
        <tr data-key="17864"><td class="text-center">1</td><td class="text-center"><img src="https://basket-12.wbbasket.ru/vol1716/part171664/171664145/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td class="articule">171664145</td><td>Контейнер для заморозки и хранения пельменей с крышкой</td><td>2023-08-03</td><td>2024-03-24</td><td class="paid">Не указана</td><td class="text-center"><a class="fa fa-floppy-o mr-auto button_paid_save" style="cursor: pointer;"></a><a class="bi bi-pencil-square mr-auto button_paid_edit" style="margin-left: 5px; cursor: pointer;"></a></td></tr>
        <tr data-key="17859"><td class="text-center">2</td><td class="text-center"><img src="https://basket-12.wbbasket.ru/vol1716/part171663/171663433/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td class="articule">171663433</td><td>Контейнер для заморозки и хранения с крышкой 1 ярус</td><td>2023-08-03</td><td>2024-11-11</td><td class="paid">Не указана</td><td class="text-center"><a class="fa fa-floppy-o mr-auto button_paid_save" style="cursor: pointer;"></a><a class="bi bi-pencil-square mr-auto button_paid_edit" style="margin-left: 5px; cursor: pointer;"></a></td></tr>
        <tr data-key="17860"><td class="text-center">3</td><td class="text-center"><img src="https://basket-15.wbbasket.ru/vol2191/part219144/219144904/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td class="articule">219144904</td><td>Мешки для робота пылесоса многоразовые</td><td>2024-03-22</td><td>2024-10-10</td><td class="paid">Не указана</td><td class="text-center"><a class="fa fa-floppy-o mr-auto button_paid_save" style="cursor: pointer;"></a><a class="bi bi-pencil-square mr-auto button_paid_edit" style="margin-left: 5px; cursor: pointer;"></a></td></tr>
        <tr data-key="17861"><td class="text-center">4</td><td class="text-center"><img src="https://basket-12.wbbasket.ru/vol1795/part179522/179522388/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td class="articule">179522388</td><td>Мешки для робота пылесоса многоразовые</td><td>2023-09-30</td><td>2024-10-10</td><td class="paid">Не указана</td><td class="text-center"><a class="fa fa-floppy-o mr-auto button_paid_save" style="cursor: pointer;"></a><a class="bi bi-pencil-square mr-auto button_paid_edit" style="margin-left: 5px; cursor: pointer;"></a></td></tr>
        <tr data-key="17862"><td class="text-center">5</td><td class="text-center"><img src="https://basket-15.wbbasket.ru/vol2191/part219144/219144905/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td class="articule">219144905</td><td>Мешки для робота пылесоса многоразовые</td><td>2024-03-22</td><td>2024-10-10</td><td class="paid">Не указана</td><td class="text-center"><a class="fa fa-floppy-o mr-auto button_paid_save" style="cursor: pointer;"></a><a class="bi bi-pencil-square mr-auto button_paid_edit" style="margin-left: 5px; cursor: pointer;"></a></td></tr>
        <tr data-key="17863"><td class="text-center">6</td><td class="text-center"><img src="https://basket-12.wbbasket.ru/vol1875/part187570/187570711/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td class="articule">187570711</td><td>Пельменница металлическая для лепки пельменей и варенников</td><td>2023-11-04</td><td>2024-03-24</td><td class="paid">Не указана</td><td class="text-center"><a class="fa fa-floppy-o mr-auto button_paid_save" style="cursor: pointer;"></a><a class="bi bi-pencil-square mr-auto button_paid_edit" style="margin-left: 5px; cursor: pointer;"></a></td></tr>
        <tr data-key="17858"><td class="text-center">7</td><td class="text-center"><img src="https://basket-12.wbbasket.ru/vol1690/part169002/169002992/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td class="articule">169002992</td><td>Пельменница форма</td><td>2023-07-13</td><td>2025-04-06</td><td class="paid">Не указана</td><td class="text-center"><a class="fa fa-floppy-o mr-auto button_paid_save" style="cursor: pointer;"></a><a class="bi bi-pencil-square mr-auto button_paid_edit" style="margin-left: 5px; cursor: pointer;"></a></td></tr>
        <tr data-key="17856"><td class="text-center">8</td><td class="text-center"><img src="https://basket-16.wbbasket.ru/vol2412/part241208/241208047/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td class="articule">241208047</td><td>Толкушка для пюре металлическая</td><td>2024-06-28</td><td>2025-08-06</td><td class="paid">Не указана</td><td class="text-center"><a class="fa fa-floppy-o mr-auto button_paid_save" style="cursor: pointer;"></a><a class="bi bi-pencil-square mr-auto button_paid_edit" style="margin-left: 5px; cursor: pointer;"></a></td></tr>
        <tr data-key="17857"><td class="text-center">9</td><td class="text-center"><img src="https://basket-24.wbbasket.ru/vol4151/part415116/415116466/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td class="articule">415116466</td><td>Толкушка для пюре металлическая</td><td>2025-05-13</td><td>2025-06-27</td><td class="paid">Не указана</td><td class="text-center"><a class="fa fa-floppy-o mr-auto button_paid_save" style="cursor: pointer;"></a><a class="bi bi-pencil-square mr-auto button_paid_edit" style="margin-left: 5px; cursor: pointer;"></a></td></tr>
        </tbody></table>
      <nav id="w1"></nav></div><!--    --></div>
    <!-- /.container-fluid -->
  </div>
</template>

<style scoped>

</style>