<script setup>

</script>

<template>
  <div class="content-header">
    <div class="container-fluid">
      <div class="row mb-2">
        <div class="col-sm-6">
          <h1 class="m-0">
            Детальная статистика                    </h1>
        </div><!-- /.col -->
        <div class="col-sm-6">
        </div><!-- /.col -->
      </div><!-- /.row -->
    </div><!-- /.container-fluid -->
  </div>
  <div class="content">
    <div class="site-login" style="margin-left: 0px">

      <form id="w0" action="/stat-detail" method="GET">
        <input type="hidden" name="StatForm[dateFrom]" value="2025-08-01">
        <input type="hidden" name="StatForm[dateTo]" value="2026-01-15">
        <input type="hidden" name="get-detail-cart--button" value="">
        <div class="form-group row field-statform-datefrom">
          <label class="col-lg-1 control-label" for="statform-datefrom">Дата с</label>
          <div class="col-lg-3"><div id="statform-datefrom-kvdate" class="input-group date"><div class="input-group-append"><span class="input-group-text kv-date-picker" title="Выбрать дату"><i class="fas fa-calendar-alt kv-dp-icon"></i></span></div> <input type="text" id="statform-datefrom" class="form-control krajee-datepicker" name="StatForm[dateFrom]" value="2025-08-01" placeholder="Выберите дату ..." data-datepicker-source="statform-datefrom-kvdate" data-datepicker-type="3" data-krajee-kvdatepicker="kvDatepicker_8596901a"></div></div>
          <div class="col-lg-8"><div class="invalid-feedback "></div></div>
        </div>
        <div class="form-group row field-statform-dateto">
          <label class="col-lg-1 control-label" for="statform-dateto">Дата по</label>
          <div class="col-lg-3"><div id="statform-dateto-kvdate" class="input-group date"><div class="input-group-append"><span class="input-group-text kv-date-picker" title="Выбрать дату"><i class="fas fa-calendar-alt kv-dp-icon"></i></span></div> <input type="text" id="statform-dateto" class="form-control krajee-datepicker" name="StatForm[dateTo]" value="2026-01-15" placeholder="Выберите дату ..." data-datepicker-source="statform-dateto-kvdate" data-datepicker-type="3" data-krajee-kvdatepicker="kvDatepicker_8596901a"></div></div>
          <div class="col-lg-8"><div class="invalid-feedback "></div></div>
        </div>

        <div class="form-group">
          <div class="col-lg-offset-1 col-lg-11">
            <button type="submit" class="btn btn-success" name="get-detail-cart--button">Показать</button>        </div>
        </div>

      </form>
      <div id="content" class="output" style="overflow-y: auto">
        <div id="w1" class="text-center" style="nax-width: 1800px;"><div class="summary">Показаны записи <b>1-4</b> из <b>4</b>.</div>
          <table class="table table-striped table-bordered"><colgroup><col>
            <col>
            <col>
            <col>
            <col>
            <col>
            <col>
            <col>
            <col>
            <col>
            <col width="50px">
            <col>
            <col>
            <col>
            <col>
            <col>
            <col></colgroup>
            <thead>
            <tr><th>#</th><th>Фото</th><th><a href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=name" data-sort="name">Название</a></th><th><a href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=nm_id" data-sort="nm_id">Артикул</a></th><th><a href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=delivery_rub" data-sort="delivery_rub">Логистика</a></th><th>Логистика (средний на единицу товара)</th><th><a href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=deduction" data-sort="deduction">Прочие удержания</a></th><th><a href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=storage_fee" data-sort="storage_fee">Хранение</a></th><th><a href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=additional_payment" data-sort="additional_payment">Доплаты</a></th><th><a href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=penalty" data-sort="penalty">Штрафы</a></th><th><a href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=rebill_logistic_cost" data-sort="rebill_logistic_cost">Возмещение издержек по перевозке</a></th><th><a class="desc" href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=ppvz_for_pay" data-sort="ppvz_for_pay">Выручка</a></th><th>Налоги %</th><th>Продажи</th><th>Возвраты</th><th>Себестоимость</th><th><a class="desc" href="/stat-detail?StatForm%5BdateFrom%5D=2025-08-01&amp;StatForm%5BdateTo%5D=2026-01-15&amp;get-detail-cart--button=&amp;sort=ppvz_for_pay" data-sort="ppvz_for_pay">~ Чистая прибыль</a></th></tr><tr id="w1-filters" class="filters"><td>&nbsp;</td><td>&nbsp;</td><td><input type="text" class="form-control" name="Stat[name]"></td><td><input type="text" class="form-control" name="Stat[nm_id]"></td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td></tr>
            </thead>
            <tfoot>
            <tr><td><b>Итого:</b></td><td>&nbsp;</td><td>&nbsp;</td><td>&nbsp;</td><td>₽ 9160</td><td>&nbsp;</td><td>₽ 8383</td><td>₽ 190</td><td>₽ 0</td><td>₽ 0</td><td>₽ 0</td><td>₽ 31161</td><td>5 %</td><td>50 шт.</td><td>0 шт.</td><td>₽ 0</td><td>₽ 12756.6</td></tr>
            </tfoot>
            <tbody>
            <tr data-key=""><td class="text-center">1</td><td class="text-center"><img src="https://basket-12.wbbasket.ru/vol1690/part169002/169002992/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td>Пельменница форма</td><td>169002992</td><td>2974.88</td><td>129.34</td><td>0.00</td><td>0.00</td><td>0.00</td><td>0.00</td><td>0</td><td>16411.15</td><td>5</td><td>23</td><td>0</td><td>0</td><td>12764.46</td></tr>
            <tr data-key=""><td class="text-center">2</td><td class="text-center"><img src="https://basket-16.wbbasket.ru/vol2412/part241208/241208047/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td>Толкушка для пюре металлическая</td><td>241208047</td><td>4232.83</td><td>184.04</td><td>0.00</td><td>0.00</td><td>0.00</td><td>0.00</td><td>0</td><td>11975.71</td><td>5</td><td>23</td><td>0</td><td>0</td><td>7355.74</td></tr>
            <tr data-key=""><td class="text-center">3</td><td class="text-center"><img src="https://basket-24.wbbasket.ru/vol4151/part415116/415116466/images/big/1.webp" alt="Фото" style="width: 150;height: 100px;"></td><td>Толкушка для пюре металлическая</td><td>415116466</td><td>1954.74</td><td>488.69</td><td>0.00</td><td>0.00</td><td>0.00</td><td>0.00</td><td>0</td><td>2775.64</td><td>5</td><td>4</td><td>0</td><td>0</td><td>779.86</td></tr>
            <tr data-key=""><td class="text-center">4</td><td class="text-center">Нет фото</td><td>Нет названия</td><td>0</td><td>0.00</td><td>0.00</td><td>8383.09</td><td>190.24</td><td>0.00</td><td>0.00</td><td>0</td><td>0</td><td>5</td><td>0</td><td>0</td><td>0</td><td>-8144.66</td></tr>
            </tbody></table>
          <nav id="w2"></nav></div>    </div>

    </div>
    <!-- /.container-fluid -->
  </div>
</template>

<style scoped>

</style>