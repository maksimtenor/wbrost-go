<script setup>

</script>

<template>
  <div class="content-header">
    <div class="container-fluid">
      <div class="row mb-2">
        <div class="col-sm-6">
          <h1 class="m-0">
            Отчет за период                    </h1>
        </div><!-- /.col -->
        <div class="col-sm-6">
        </div><!-- /.col -->
      </div><!-- /.row -->
    </div><!-- /.container-fluid -->
  </div>
  <div class="content">
    <div class="site-login">
    <span style="color: red; font-weight: bold;"> Внимание! Из-за высокой нагрузки, временно добавлено ограничение в 1 месяц! <br>
    Только на загрузку новой... Просмотр уже имеющейся статистики, остается таким же!</span>
      <form id="w0" class="form-horizontal" action="/index" method="post">
        <input type="hidden" name="_csrf" value="92UmL_hhhZfn9iBLuMNPBm7rR23zzRGXqZt54FdZlS6gFQtVqw_GyIa_EAGKkBpsJKMKA6aefsXP_UupZBL5bQ==">
        <!--    --><!---->
        <!--    -->
        <div class="form-group field-statform-datefrom">
          <label class="col-lg-1 control-label" for="statform-datefrom">Дата с</label>
          <div class="col-lg-3"><div id="statform-datefrom-kvdate" class="input-group date"><div class="input-group-append"><span class="input-group-text kv-date-picker" title="Выбрать дату"><i class="fas fa-calendar-alt kv-dp-icon"></i></span></div> <input type="text" id="statform-datefrom" class="form-control krajee-datepicker" name="StatForm[dateFrom]" value="2026-01-08" placeholder="Выберите дату ..." data-datepicker-source="statform-datefrom-kvdate" data-datepicker-type="3" data-krajee-kvdatepicker="kvDatepicker_4628cdd9"></div></div>
          <div class="col-lg-8"><p class="help-block help-block-error "></p></div>
        </div>
        <div class="form-group field-statform-dateto">
          <label class="col-lg-1 control-label" for="statform-dateto">Дата по</label>
          <div class="col-lg-3"><div id="statform-dateto-kvdate" class="input-group date"><div class="input-group-append"><span class="input-group-text kv-date-picker" title="Выбрать дату"><i class="fas fa-calendar-alt kv-dp-icon"></i></span></div> <input type="text" id="statform-dateto" class="form-control krajee-datepicker" name="StatForm[dateTo]" value="2026-01-15" placeholder="Выберите дату ..." data-datepicker-source="statform-dateto-kvdate" data-datepicker-type="3" data-krajee-kvdatepicker="kvDatepicker_4628cdd9"></div></div>
          <div class="col-lg-8"><p class="help-block help-block-error "></p></div>
        </div>
        <!--    -->
        <div class="form-group">
          <div class="col-lg-offset-1 col-lg-11">
            <button type="submit" class="btn btn-success" name="get-detail-cart--button">Запросить</button>        </div>
        </div>

      </form>
      <div id="content" class="output" style="overflow-y: auto">
        <div id="w1" class="text-center" style="nax-width: 1800px;"><div class="summary">Показаны записи <b>1-10</b> из <b>35</b>.</div>
          <table class="table table-striped table-bordered"><colgroup><col style="max-width:100px">
            <col>
            <col>
            <col>
            <col style="max-width:300px"></colgroup>
            <thead>
            <tr><th>Статус</th><th>Диапазон</th><th>Дата создания</th><th>Дата последнего изменения</th><th>Комментарий к ошибке</th></tr>
            </thead>
            <tbody>
            <tr data-key="51"><td class="text-center"><span style="color: green">Готово</span></td><td class="text-center">2025-09-02 - 2025-09-09</td><td class="text-center">2025-09-09 16:01:34</td><td class="text-center">2025-09-09 13:02:02</td><td class="text-center"><span style="color: orange"></span></td></tr>
            <tr data-key="50"><td class="text-center"><span style="color: green">Готово</span></td><td class="text-center">2025-08-25 - 2025-09-01</td><td class="text-center">2025-09-01 13:53:28</td><td class="text-center">2025-09-01 10:54:02</td><td class="text-center"><span style="color: orange"></span></td></tr>
            <tr data-key="49"><td class="text-center"><span style="color: red">Ошибка</span></td><td class="text-center">2025-08-04 - 2025-08-24</td><td class="text-center">2025-08-24 21:49:22</td><td class="text-center">2025-08-24 18:50:02</td><td class="text-center"><span style="color: red">Ошибка обработки</span></td></tr>
            <tr data-key="46"><td class="text-center"><span style="color: green">Готово</span></td><td class="text-center">2025-07-16 - 2025-07-23</td><td class="text-center">2025-07-23 09:54:24</td><td class="text-center">2025-07-23 06:55:03</td><td class="text-center"><span style="color: orange"></span></td></tr>
            <tr data-key="45"><td class="text-center"><span style="color: green">Готово</span></td><td class="text-center">2025-07-07 - 2025-07-14</td><td class="text-center">2025-07-14 18:26:49</td><td class="text-center">2025-07-14 15:27:03</td><td class="text-center"><span style="color: orange"></span></td></tr>
            <tr data-key="44"><td class="text-center"><span style="color: green">Готово</span></td><td class="text-center">2025-06-10 - 2025-07-10</td><td class="text-center">2025-07-10 23:16:22</td><td class="text-center">2025-07-10 20:17:07</td><td class="text-center"><span style="color: orange">Отчет успешно получен и загружен в систему!
                                <br> Получено 2312 строк !
                                <br> Сохранено 2310 строк !
                                <br> Не сохранилось 2 строк ! </span></td></tr>
            <tr data-key="43"><td class="text-center"><span style="color: green">Готово</span></td><td class="text-center">2025-06-02 - 2025-06-09</td><td class="text-center">2025-06-16 14:46:48</td><td class="text-center">2025-06-16 11:47:04</td><td class="text-center"><span style="color: orange"></span></td></tr>
            <tr data-key="42"><td class="text-center"><span style="color: red">Ошибка</span></td><td class="text-center">2025-06-09 - 2025-06-16</td><td class="text-center">2025-06-16 14:45:48</td><td class="text-center">2025-06-16 11:46:02</td><td class="text-center"><span style="color: red">Нет данных</span></td></tr>
            <tr data-key="40"><td class="text-center"><span style="color: green">Готово</span></td><td class="text-center">2025-05-25 - 2025-06-01</td><td class="text-center">2025-06-01 18:22:59</td><td class="text-center">2025-06-01 15:23:02</td><td class="text-center"><span style="color: orange">Что-то пошло не так при сохранении в БД
                                <br> Получено 106 строк !
                                <br> Сохранено 0 строк !
                                <br> Не сохранилось 106 строк (Возможно уже есть в системе)!
                            </span></td></tr>
            <tr data-key="39"><td class="text-center"><span style="color: green">Готово</span></td><td class="text-center">2025-05-01 - 2025-05-30</td><td class="text-center">2025-05-30 09:33:46</td><td class="text-center">2025-05-30 06:34:06</td><td class="text-center"><span style="color: orange">Отчет успешно получен и загружен в систему!
                                <br> Получено 2105 строк !
                                <br> Сохранено 2021 строк !
                                <br> Не сохранилось 84 строк ! </span></td></tr>
            </tbody></table>
          <nav id="w2"><ul class="pagination"><li class="page-item prev disabled"><a class="page-link" href="/index?page=1&amp;per-page=10" data-page="0" tabindex="-1"><span aria-hidden="true">«</span>
            <span class="sr-only">Previous</span></a></li>
            <li class="page-item active"><a class="page-link" href="/index?page=1&amp;per-page=10" data-page="0">1</a></li>
            <li class="page-item"><a class="page-link" href="/index?page=2&amp;per-page=10" data-page="1">2</a></li>
            <li class="page-item"><a class="page-link" href="/index?page=3&amp;per-page=10" data-page="2">3</a></li>
            <li class="page-item"><a class="page-link" href="/index?page=4&amp;per-page=10" data-page="3">4</a></li>
            <li class="page-item next"><a class="page-link" href="/index?page=2&amp;per-page=10" data-page="1"><span aria-hidden="true">»</span>
              <span class="sr-only">Next</span></a></li></ul></nav></div>    </div>

    </div>


    <!-- /.container-fluid -->
  </div>
</template>

<style scoped>

</style>